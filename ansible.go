package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"EasyRvB/service"

	"github.com/apenella/go-ansible/v2/pkg/execute"
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

func RunPlaybook(role string, ip net.IP) error {
	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Become:    true,
		Inventory: ip.String() + ",",
		Tags:      role,
	}

	type roleTemplate struct {
		Role string
	}

	tmpl, err := template.New("playbook.tmpl").ParseGlob("ansible/playbook.tmpl")
	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create("ansible/playbook.yaml")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, roleTemplate{role})
	if err != nil {
		panic(err)
	}

	playbookCmd := playbook.NewAnsiblePlaybookCmd(
		playbook.WithPlaybooks("ansible/playbook.yaml"),
		playbook.WithPlaybookOptions(ansiblePlaybookOptions),
	)

	exec := execute.NewDefaultExecute(
		execute.WithCmd(playbookCmd),
	)

	err = exec.Execute(context.Background())
	if err != nil {
		return err
	}

	return nil
}
