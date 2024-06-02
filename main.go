package main

import (
	"EasyRvB/service"
	"fmt"
	"os"
)

func main() {
	roles, err := getAnsibleRoles("./ansible")
	if err != nil {
		fmt.Println(err)
	}

	for _, role := range roles {
		service := CreateService(role)
		printServiceDetails(service)
	}
}

func CreateService(role string) service.Service {
	roleType, err := getRoleType(role)
	if err != nil {
		fmt.Println(err)
	}

	switch roleType {
	case "kubernetes":
		service := service.KubernetesConfig{}
		err := service.ReadConfig(role + string(os.PathSeparator) + roleType + ".toml")
		if err != nil {
			fmt.Println(err)
		}
		return &service
	}
	return nil
}

func printServiceDetails(s service.Service) {
	switch s := s.(type) {
	case *service.KubernetesConfig:
		fmt.Printf("Service Name: %s\n", s.Name)
	}
}
