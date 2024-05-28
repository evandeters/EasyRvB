package main

import (
    "fmt"
    "os"
    "path/filepath"
)

type Config struct {
    OsSupport map[string][]string
}

func getConfigFromFile(path string) Config {
    var config Config
    config.OsSupport = make(map[string][]string)

    file, err := os.Open(path)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return config
    }
    defer file.Close()

    
