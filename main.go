package main

import (
	"EasyRvB/service"
	"errors"
	"fmt"
	"os"
)

func main() {
	roles, err := getAnsibleRoles("./ansible")
	if err != nil {
		fmt.Println(err)
	}

	var ServiceConfigs []*service.ServiceConfig

	for _, role := range roles {
		service, err := CreateService(role)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ServiceConfigs = append(ServiceConfigs, service)
	}
	for i, svc := range ServiceConfigs {
		fmt.Println(ServiceConfigs[i].Name)
		fmt.Println(ServiceConfigs[i].Port)

		fmt.Println(svc)
	}
}

func CreateService(role string) (*service.ServiceConfig, error) {
	roleType, err := getRoleType(role)
	if err != nil {
		fmt.Println(err)
	}
	configPath := role + string(os.PathSeparator) + roleType + ".toml"

	if _, err := os.Stat(configPath); err != nil {
		return nil, err
	}

	service := service.ServiceConfig{}
	err = service.ReadConfig(configPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Role Type Unknown: %v. Error: %v", roleType, err))
	}

	return &service, nil
}
