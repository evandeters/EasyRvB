package service

import (
	"fmt"
	"os"
	"text/template"
)

type HTTPConfig struct {
	Host string `toml:"host"`
	RequireDB  bool `toml:"requiredb"`
    DatabaseIP string `toml:"databaseip"`
    VarTmplPath string `toml:"vartmplpath"`
    VarYamlPath string `toml:"varyamlpath"`
}

func (h HTTPConfig) SetConfig(role string) error {
    httpData := struct {
        Host string
        DatabaseIP string
        Port int
    }{
        Host: h.Host,
        DatabaseIP: h.DatabaseIP,
        Port: 80,
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

    err = tmpl.Execute(outputFile, httpData)
    if err != nil {
        panic(err)
    }
    return nil
}

func (h *HTTPConfig) SetDatabaseIP(ip string) {
    h.DatabaseIP = ip
}
