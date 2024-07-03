package main

import (
    "bytes" 
    "context"
    "math/rand"
    "errors"
    "fmt"
    "io"
    "net"
    "os"
    "path/filepath"
    "strings"

    "EasyRvB/host"
    "EasyRvB/service"

    "github.com/apenella/go-ansible/pkg/stdoutcallback/results"
    "github.com/apenella/go-ansible/v2/pkg/execute"
    "github.com/apenella/go-ansible/v2/pkg/execute/stdoutcallback"
    "github.com/apenella/go-ansible/v2/pkg/playbook"
)

func getAnsibleRoles(path string) ([]string, error) {
    var roles []string
    error := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        pathArr := strings.Split(filePath, string(os.PathSeparator))

        if pathArr[len(pathArr)-2] == "roles" && !strings.HasPrefix(pathArr[len(pathArr)-1], ".") {
            roles = append(roles, filePath)
        }

        return nil
    })

    if error != nil {
        return nil, error
    }

    return roles, nil
}

func GetRoleType(path string) (string, error) {
    var fileData string
    error := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        pathArr := strings.Split(filePath, string(os.PathSeparator))

        if strings.HasSuffix(pathArr[len(pathArr)-1], ".toml") {
            fileData = pathArr[len(pathArr)-1]
        }

        return nil
    })

    if error != nil {
        return "", error
    }

    return strings.Split(fileData, ".")[0], nil
}

func ServiceFromRole(role, roleType string) (*service.ServiceConfig, error) {
    fmt.Printf("Role: %v, RoleType: %v\n", role, roleType)

    configPath := role + string(os.PathSeparator) + roleType + ".toml"

    if _, err := os.Stat(configPath); err != nil {
        return nil, err
    }

    service := service.ServiceConfig{}
    err := service.ReadConfig(configPath)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Role Type Unknown: %v. Error: %v", role, err))
    }

    return &service, nil
}

func RunPlaybook(role string, target *host.Host, user string) error {
    ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
        Become:    true,
        Inventory: target.NattedIp.String() + ",",
        Tags:      role,
        User:      user,

    }


    type roleTemplate struct {
        Role string
    }

    FillTemplate("ansible/templates/playbook.tmpl", roleTemplate{Role: role})

    playbookCmd := playbook.NewAnsiblePlaybookCmd(
        playbook.WithPlaybooks("ansible/playbook.yaml"),
        playbook.WithPlaybookOptions(ansiblePlaybookOptions),
    )

    exec := execute.NewDefaultExecute(
        execute.WithCmd(playbookCmd),
    )

    err := exec.Execute(context.Background())
    if err != nil { return err }

    return nil
}

func CreateVM(templateName string, vmName string) net.IP {
    var boxIP string
    for {
        ipExists := false
        fourthOctet := rand.Intn(250) + 4
        boxIP = fmt.Sprintf("192.168.1.%v", fourthOctet)
        for _, box := range CurrentHosts {
            if box.Ip.String() == boxIP {
                ipExists = true
                break
            }
        }
        if !ipExists {
            break
        }
    }

    vmData := struct{
        TemplateName string
        VCenterServer string
        VCenterUser string
        VCenterPassword string
        ResourcePool string
        Datacenter string
        Cluster string
        VMFolder string
        VMName string
        Datastore string
        PortGroup string
        IPAddress string
    }{
        TemplateName: templateName,
        VCenterServer: ConfigMap.VCenterServer,
        VCenterUser: ConfigMap.VCenterUsername,
        VCenterPassword: ConfigMap.VCenterPassword,
        ResourcePool: ConfigMap.ResourcePool,
        Datacenter: ConfigMap.Datacenter,
        Cluster: ConfigMap.Cluster,
        VMFolder: ConfigMap.VMFolder,
        Datastore: ConfigMap.Datastore,
        PortGroup: ConfigMap.PortGroup,
        VMName: vmName,
        IPAddress: boxIP,
    }

    FillTemplate("ansible/templates/vm.tmpl", vmData)
    ReplaceBrackets("ansible/vm.yaml")

    playbookCmd := playbook.NewAnsiblePlaybookCmd(
        playbook.WithPlaybooks("ansible/vm.yaml"),
    )

    buff := new(bytes.Buffer)

    exec := stdoutcallback.NewJSONStdoutCallbackExecute(execute.NewDefaultExecute(
        execute.WithCmd(playbookCmd),
        execute.WithWrite(io.Writer(buff)),
        ),
    )

    err := exec.Execute(context.Background())
    if err != nil { 
        panic(err)
    }

    res, err := results.ParseJSONResultsStream(buff)
    if err != nil {
        panic(err)
    }

    ipAdd := res.Plays[0].Tasks[1].Hosts["localhost"].AnsibleFacts["ip_dict"].(map[string]interface{})["ip_address"]
    if ipAdd == nil {
        panic("IP Address not found")
    }

    return net.ParseIP(ipAdd.(string))
}

func CreateRouter(octet int) *host.Host {
    vmName := fmt.Sprintf("router-%v", octet)
    routerData := struct{
        octet int
    }{
        octet: octet,
    }

    FillTemplate("ansible/templates/router.tmpl", routerData)

    ip := CreateVM(ConfigMap.RouterTemplate, vmName)

    router := host.Host{
        Hostname: vmName,
        Ip: ip,
    }

    return &router
}
