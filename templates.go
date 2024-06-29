package main

import (
    "os"
    "strings"
    "text/template"
)

func ReplaceBrackets(filePath string) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        panic(err)
    }

    lines := strings.Split(string(data), "\n")
    for i, line := range lines {
        if strings.Contains(line, "[[") && strings.Contains(line, "]]") {
            lines[i] = strings.Replace(line, "[[", "{{", -1)
            lines[i] = strings.Replace(lines[i], "]]", "}}", -1)
        }
    }

    output := strings.Join(lines, "\n")
    err = os.WriteFile(filePath, []byte(output), 0644)
    if err != nil {
        panic(err)
    }
}

func FillTemplate(filePath string, data interface{}) {
    templateName := strings.Split(filePath, string(os.PathSeparator))

    tmpl, err := template.New(templateName[len(templateName) - 1]).ParseGlob(filePath)
    if err != nil {
        panic(err)
    }

    fileName := strings.Split(templateName[len(templateName) - 1], ".")[0]

    outputFile, err := os.Create("ansible/" + fileName + ".yaml")
    if err != nil {
        panic(err)
    }
    defer outputFile.Close()

    err = tmpl.Execute(outputFile, data)
    if err != nil {
        panic(err)
    }
}

