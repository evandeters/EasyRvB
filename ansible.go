package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"EasyRvB/host"
	"EasyRvB/service"

	"github.com/apenella/go-ansible/v2/pkg/execute"
	"github.com/apenella/go-ansible/v2/pkg/playbook"
)

func GetAnsibleRoles(path string) ([]string, error) {
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

func ServiceFromRole(role string) (service.ServiceConfig, error) {
    files, err := os.ReadDir(role)
    if err != nil {
        return service.ServiceConfig{}, err
    }
    var configPath string
    for _, f := range files {
        if !f.IsDir() && strings.HasSuffix(f.Name(), ".toml") {
             configPath = filepath.Join(role, f.Name())
        }
    }

    if _, err := os.Stat(configPath); err != nil {
        return service.ServiceConfig{}, err
    }

    svc := service.ServiceConfig{}
    err = svc.ReadConfig(configPath)
    if err != nil {
        return service.ServiceConfig{}, errors.New(fmt.Sprintf("Role Type Unknown: %v. Error: %v", role, err))
    }

    return svc, nil
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

    FillTemplate("ansible/templates/playbook.tmpl", "ansible/playbook.yaml", roleTemplate{Role: role})

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

