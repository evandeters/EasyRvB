package service

import (
    "fmt"
)

type HTTPConfig struct {
    Host string `json:"host"`
    Cms string `json:"cms"`
}

func (c *HTTPConfig) ReadConfig() error {
    fmt.Println("Reading config")
    
    fmt.Println("Host: ", c.Host)
    fmt.Println("CMS: ", c.Cms)

    return nil
}
