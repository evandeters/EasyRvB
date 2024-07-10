package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"

	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"github.com/apenella/go-ansible/v2/pkg/execute"
	"github.com/apenella/go-ansible/v2/pkg/execute/stdoutcallback"
	"github.com/apenella/go-ansible/v2/pkg/playbook"
)

func CreateVM(vm Host) error {
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
        TemplateName: vm.VMTemplate,
        VCenterServer: ConfigMap.VCenterServer,
        VCenterUser: ConfigMap.VCenterUsername,
        VCenterPassword: ConfigMap.VCenterPassword,
        ResourcePool: ConfigMap.ResourcePool,
        Datacenter: ConfigMap.Datacenter,
        Cluster: ConfigMap.Cluster,
        VMFolder: ConfigMap.VMFolder,
        Datastore: ConfigMap.Datastore,
        PortGroup: ConfigMap.PortGroup,
        VMName: vm.Hostname,
        IPAddress: vm.Ip.String(),
    }

    FillTemplate("ansible/templates/vm.tmpl", "ansible/vm.yaml", vmData)
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

    _, err = results.ParseJSONResultsStream(buff)
    if err != nil {
        return err
    }

    return nil
}

func CreateRouter(templateName string) *Host {
    vmName := fmt.Sprintf("router-%v", ThirdOctet)

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
        VMName: vmName,
    }

    FillTemplate("ansible/templates/router.tmpl", "ansible/router.yaml", vmData)
    ReplaceBrackets("ansible/router.yaml")

    playbookCmd := playbook.NewAnsiblePlaybookCmd(
        playbook.WithPlaybooks("ansible/router.yaml"),
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

    fmt.Println(res)

    ip := res.Plays[0].Tasks[2].Hosts["localhost"].AnsibleFacts["ip_dict"].(map[string]interface{})["ip_address"]
    if ip == nil {
        panic("IP Address not found")
    }

    ipAdd := net.ParseIP(ip.(string))

    newHost := Host{
        Hostname: vmName,
        Ip: ipAdd,
    }

    return &newHost
}

func ConfigRouter(router *Host, octet int) error {
    routerData := struct{
        Octet int
    }{
        Octet: octet,
    }

    inventoryData := struct{
        RouterIp string
    }{
        RouterIp: router.Ip.String(),
    }

    ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
        Become:    true,
        Inventory: "ansible/inventory/inventory.ini",
        Tags:      "router",
    }

    FillTemplate("ansible/roles/router/vars/default.tmpl", "ansible/roles/router/vars/default.yaml", routerData)
    FillTemplate("ansible/templates/inventory.tmpl", "ansible/inventory/inventory.ini", inventoryData)

    type roleTemplate struct {
        Role string
    }

    FillTemplate("ansible/templates/playbook.tmpl", "ansible/playbook.yaml", roleTemplate{Role: "router"})

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
