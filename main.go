package main

import (
	"EasyRvB/host"
	"EasyRvB/service"
	"fmt"
)

var ServiceConfigs map[string]service.ServiceConfig
var ServiceNames []string
var ConfigMap Config
var CurrentHosts []*host.Host
var ThirdOctet int

func init() {

    ConfigMap = Config{}
    ReadConfig(
    	&ConfigMap,
    	"config.toml",
    )
    fmt.Println("VMTemplates", ConfigMap.VMTemplates)
    ThirdOctet = GetAvailableOctet()
    ServiceConfigs = make(map[string]service.ServiceConfig)

    roles, err := GetAnsibleRoles("./ansible")
    if err != nil {
        fmt.Println(err)
    }

    for _, role := range roles {
        svc, err := ServiceFromRole(role)
        if err != nil {
            fmt.Println(err)
        }
        if svc.Name != "" {
            ServiceConfigs[svc.Name] = svc
        }
    }

    for key, svc := range ServiceConfigs {
        ServiceNames = append(ServiceNames, key)
        if (svc.Http != service.HTTPConfig{}) {
            fmt.Println("HTTP Service found:", svc.Name)
        }

        if (svc.Database != service.DatabaseConfig{}) {
            fmt.Println("Database Service found:", svc.Name)
        }

        if (svc.Kubernetes != service.KubernetesConfig{}) {
            fmt.Println("Kubernetes Service found:", svc.Name)
        }
    }

}

func main() {
    /**router := CreateRouter("Cisco CSR1kv Blank")
    CurrentHosts = append(CurrentHosts, router)

    fmt.Println("Router created:", router.Hostname, router.Ip)
    fmt.Println("Configuring router...")
    err := ConfigRouter(router, ThirdOctet)
    if err != nil {
        fmt.Println(err)
    }*/

    GenerateHosts(2)
    for _, host := range CurrentHosts {
        fmt.Println("Host generated:", host.Hostname, host.Ip, host.NattedIp, host.VMTemplate)
    }

    AddServicesToHosts()
    for _, host := range CurrentHosts {
        host.GetServices()
    }

    for _, host := range CurrentHosts {
        vm := CreateVM(host.VMTemplate, host.Hostname)
        for _, svc := range host.Services {
            err := RunPlaybook(svc.ConfigMap.Name, vm, svc.ConfigMap.User)
            if err != nil {
                fmt.Println(err)
            }
        }
    }
    

    /*
    vm := CreateVM("Ubuntu 22.04 Blank", "apache-vm")
    CurrentHosts = append(CurrentHosts, vm)
    fmt.Println("VM created:", vm.Hostname, vm.Ip)
    fmt.Println("Configuring VM...")
    err = RunPlaybook("apache", vm, "root")
    if err != nil {
        fmt.Println(err)
    }

    vm := CreateVM("Ubuntu 22.04 Blank", "mysql-vm")
    CurrentHosts = append(CurrentHosts, vm)
    vm.AddService(service.ServiceInstance{ID: uuid.New(), ConfigMap: ServiceConfigs["wordpress_mysql"]})
    fmt.Println("VM created:", vm.Hostname, vm.Ip)
    fmt.Println("Configuring VM...")
    err := RunPlaybook("wordpress-db", vm, "root")
    if err != nil {
       fmt.Println(err)
    }

    vm = CreateVM("Ubuntu 22.04 Blank", "wordpress-vm")
    CurrentHosts = append(CurrentHosts, vm)
    vm.AddService(service.ServiceInstance{ID: uuid.New(), ConfigMap: ServiceConfigs["wordpress_http"]})
    for _, svc := range vm.Services {
        if svc.ConfigMap.Http.RequireDB {
            for _, host := range CurrentHosts {
                for _, service := range host.Services {
                    if strings.Contains(service.ConfigMap.Name, "mysql") {
                        FillTemplate("ansible/roles/wordpress-web/vars/default.tmpl", "ansible/roles/wordpress-web/vars/default.yaml", map[string]string{"Ip": host.Ip.String()})
                    }
                }
            }
        }
    }
    fmt.Println("VM created:", vm.Hostname, vm.Ip)
    fmt.Println("Configuring VM...")
    err = RunPlaybook("wordpress-web", vm, "root")
    vm.GetServices()
    */

}

func GetAvailableOctet() int {
    // To implement
    return 245
}
