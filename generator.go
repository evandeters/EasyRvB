package main

import (
	"EasyRvB/host"
	"EasyRvB/service"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func GenerateHosts(numHosts int) {
    for i := 0; i < numHosts; i++ {
        random := rand.Intn(len(ConfigMap.VMTemplates))
        ip := GenerateIP("192.168.1")
        nattedIP := "172.16." + strconv.Itoa(ThirdOctet) + "." + strings.Split(ip.String(), ".")[3]
        newHost := host.Host{
            Hostname: gofakeit.CarMaker(),
            Ip: ip,
            NattedIp: net.ParseIP(nattedIP),
            VMTemplate: ConfigMap.VMTemplates[random],
        }
        CurrentHosts = append(CurrentHosts, &newHost)
    }
}

func AddServicesToHosts() {
    for _, vm := range CurrentHosts {
        for i := 0; i < 1; i++ {
            svc := service.ServiceInstance{
                ID: uuid.New(),
                ConfigMap: ServiceConfigs[ServiceNames[rand.Intn(len(ServiceNames))]],
            }
            for _, os := range svc.ConfigMap.SupportedOS {
                fmt.Println("Checking if", vm.VMTemplate, "supports", os)
                if strings.Contains(strings.ToLower(vm.VMTemplate), strings.ToLower(os)) {
                    vm.Services = append(vm.Services, &svc)
                }
            }
        }
    }
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
