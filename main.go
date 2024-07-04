package main

import (
	"EasyRvB/host"
	"EasyRvB/service"
	"fmt"
)

var ServiceConfigs map[string]*service.ServiceConfig
var ConfigMap Config
var CurrentHosts []*host.Host
var ThirdOctet int

func init() {

    ConfigMap = Config{}
    ReadConfig(&ConfigMap, "config.toml")
    ThirdOctet = getAvailableOctet()
	ServiceConfigs = make(map[string]*service.ServiceConfig)

	roles, err := getAnsibleRoles("./ansible")
	if err != nil {
		fmt.Println(err)
	}

	for _, role := range roles {
		roleType, err := GetRoleType(role)
		if err != nil {
			panic(err)
		}

		service, err := ServiceFromRole(role, roleType)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ServiceConfigs[service.Name] = service
	}

	for _, svc := range ServiceConfigs {
		if (svc.Http != service.HTTPConfig{}) {
			// do stuff
		}

		if (svc.Kubernetes != service.KubernetesConfig{}) {
		}
	}

}

func main() {
    router := CreateRouter("Cisco CSR1kv Blank")
    CurrentHosts = append(CurrentHosts, router)

    fmt.Println("Router created:", router.Hostname, router.Ip)
    fmt.Println("Configuring router...")
    err := ConfigRouter(router, ThirdOctet)
    if err != nil {
        fmt.Println(err)
    }

    vm := CreateVM("Ubuntu 22.04 Blank", "apache-vm")
    CurrentHosts = append(CurrentHosts, vm)
    fmt.Println("VM created:", vm.Hostname, vm.Ip)
    fmt.Println("Configuring VM...")
    err = RunPlaybook("apache", vm, "root")
    if err != nil {
        fmt.Println(err)
    }
}

func getAvailableOctet() int {
    // To implement
    return 245
}
