package main

import (
	"EasyRvB/host"
	"EasyRvB/service"
	"fmt"
    "time"
	"net"
)

var ServiceConfigs map[string]*service.ServiceConfig
var ConfigMap Config
var CurrentHosts []*host.Host
var ThirdOctet int

func init() {

    ConfigMap = Config{}
    ReadConfig(&ConfigMap, "config.toml")
    fmt.Println(ConfigMap)

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
		fmt.Println(service.Name)
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
    ip := CreateVM("Ubuntu 22.04 Blank", "TestVM4")
    nattedIP := net.IPv4(172, 16, byte(ThirdOctet), ip.To4()[3])
    testHost := host.NewHost("test-web", ip, nattedIP, "Ubuntu22")
    CurrentHosts = append(CurrentHosts, testHost)

    time.Sleep(30 * time.Second)

	err := RunPlaybook("apache", testHost, "root")
	if err != nil {
		fmt.Println(err)
	} else {
		testHost.AddService(*service.NewServiceInstance(ServiceConfigs["apache_http"]))
	}

	testHost.GetServices()
}

func getAvailableOctet() int {
    // To implement
    return 245
}
