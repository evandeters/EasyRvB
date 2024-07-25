package service

import (
	"fmt"
	"os"
	"text/template"
)

type DnsConfig struct {
    Domain string `toml:"domain"`
    VarTmplPath string `toml:"vartmplpath"`
    VarYamlPath string `toml:"varyamlpath"`
}

func (h DnsConfig) SetConfig(role string) error {
    dnsData := struct {
        Domain string
    }{
        Domain: h.Domain,
    }
    templatePath := "ansible/roles/" + role + "/" + h.VarTmplPath
    resultPath := "ansible/roles/" + role + "/" + h.VarYamlPath

    fmt.Println("Template Path:", templatePath)
    fmt.Println("Result Path:", resultPath)
    tmpl, err := template.New("default.tmpl").ParseGlob(templatePath)
    if err != nil {
        panic(err)
    }

    outputFile, err := os.Create(resultPath)
    if err != nil {
        panic(err)
    }
    defer outputFile.Close()

    err = tmpl.Execute(outputFile, dnsData)
    if err != nil {
        panic(err)
    }
    return nil
}
