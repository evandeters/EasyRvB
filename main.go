package main

import (
	"EasyRvB/host"
	"EasyRvB/service"
	"fmt"
)

var ServiceConfigs map[string]*service.ServiceConfig
var ConfigMap Config

func init() {

    ConfigMap = Config{}
    ReadConfig(&ConfigMap, "config.toml")
    fmt.Println(ConfigMap)
    
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
    ip := CreateVM("Ubuntu 22.04 Blank", "TestVM3")
    testHost := host.NewHost("test-web", ip, "Ubuntu22")
	err := RunPlaybook("apache", testHost)
	if err != nil {
		fmt.Println(err)
	} else {
		testHost.AddService(*service.NewServiceInstance(ServiceConfigs["apache_http"]))
	}

	testHost.GetServices()
}
