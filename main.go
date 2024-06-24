package main

import (
	"EasyRvB/host"
	"EasyRvB/service"
	"fmt"
	"net"
)

var ServiceConfigs map[string]*service.ServiceConfig

func init() {
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
	target := net.IPv4(192, 168, 68, 30)
	testHost := host.NewHost("test-web", target, "Ubuntu22")
	fmt.Println(target)
	err := RunPlaybook("kube-master-node", target)
	if err != nil {
		fmt.Println(err)
	} else {
		testHost.AddService(*service.NewServiceInstance(ServiceConfigs["apache_http"]))
	}

	testHost.GetServices()
}
