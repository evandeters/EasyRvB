package main

import (
    "fmt"
    "math/rand"
    "net"
    "strconv"
    "strings"

    "github.com/brianvoe/gofakeit/v7"
)

func GenerateHosts(numHosts int) {
    for i := 0; i < numHosts; i++ {
        random := rand.Intn(len(ConfigMap.VMTemplates))
        ip := GenerateIP("192.168.1")
        nattedIP := "172.16." + strconv.Itoa(ThirdOctet) + "." + strings.Split(ip.String(), ".")[3]
        newHost := Host{
            Hostname: gofakeit.CarMaker(),
            Ip: ip,
            NattedIp: net.ParseIP(nattedIP),
            VMTemplate: ConfigMap.VMTemplates[random],
        }
        CurrentHosts = append(CurrentHosts, &newHost)
        err := AddHost(newHost)
        if err != nil {
            fmt.Println("Error:", err)
        }
    }
}

func AddServicesToHosts() {
    hosts, err := GetHosts()
    if err != nil {
        fmt.Println(err)
    }
    for _, host := range hosts {
        numServices := rand.Intn(MaxServices) + 1
        for i := 0; i < numServices; i++ {
            svcName := ServiceNames[rand.Intn(len(ServiceNames))]
            AddService(host, svcName)
            HandleDependency(svcName)
        }
    }
}

func HandleDependency(svcName string) {
    for _, dep := range ServiceConfigs[svcName].Dependency {
        ips := GetDependencyIPs(dep)
        if ips == nil {
            randHost := *CurrentHosts[rand.Intn(len(CurrentHosts))]
            AddService(randHost, dep)
            ips = append(ips, randHost.Ip.String())
        }
        switch ServiceConfigs[svcName].Type {
        case "http":
            ServiceConfigs[svcName].Http.SetDatabaseIP(ips[rand.Intn(len(ips))])
            ServiceConfigs[svcName].Http.SetConfig(svcName)
        }
    }
}

func GetDependencyIPs(service string) []string {
    var dependencyIPs []string
    hosts, err := GetHostsWithService(service)
    if err != nil {
        fmt.Println(err)
    }
    for _, host := range hosts {
        dependencyIPs = append(dependencyIPs, host.Ip.String())
    }
    return dependencyIPs
}

func GenerateIP(network string) net.IP {
    var boxIP string
    for {
        ipExists := false
        fourthOctet := rand.Intn(250) + 4
        boxIP = fmt.Sprintf("%v.%v", network, fourthOctet)
        for _, box := range CurrentHosts {
            if box.Ip.String() == boxIP {
                ipExists = true
                break
            }
        }
        if !ipExists {
            break
        }
    }
    return net.ParseIP(boxIP)
}
