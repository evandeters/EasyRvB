package main

import (
    "fmt"
    "net"

    "EasyRvB/service"
    "EasyRvB/host"
)

func main() {
    webServer := host.NewHost("webServer", net.ParseIP("192.168.1.50"), "Linux")
    webServer.AddService(*service.NewService("webServer", 80, "TCP", nil, nil))
    webConfig, err := service.NewServiceConfig("HTTP", []byte(`{"host": "localhost", "cms": "wordpress"}`))
    webServer.Services[0].Config = webConfig
    if err != nil {
        fmt.Println(err)
        return
    }
    webServer.Services[0].Config.ReadConfig()

    webServer.Println()

}
