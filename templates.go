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

func FillTemplate(tmplPath, resultPath string, data interface{}) {
    templateName := strings.Split(tmplPath, string(os.PathSeparator))

    tmpl, err := template.New(templateName[len(templateName) - 1]).ParseGlob(tmplPath)
    if err != nil {
        panic(err)
    }


    outputFile, err := os.Create(resultPath)
    if err != nil {
        panic(err)
    }
    defer outputFile.Close()

    err = tmpl.Execute(outputFile, data)
    if err != nil {
        panic(err)
    }
}

