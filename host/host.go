package host

import (
	"EasyRvB/service"
	"fmt"
	"net"
)

type Host struct {
	Hostname        string
	Ip              net.IP
    NattedIp        net.IP
	VMTemplate        string
	Services        []*service.ServiceInstance
}

func NewHost(hostname string, ip, nattedIP net.IP, template string) *Host {
	return &Host{
		Hostname:        hostname,
		Ip:              ip,
        NattedIp:        nattedIP,
		VMTemplate:      template,
		Services:        []*service.ServiceInstance{},
	}
}

func (h *Host) AddService(s service.ServiceInstance) {
	h.Services = append(h.Services, &s)
}

func (h *Host) RemoveService(s service.ServiceInstance) {
	for i, service := range h.Services {
		if service.ID == s.ID {
			h.Services = append(h.Services[:i], h.Services[i+1:]...)
			break
		}
	}
}

func (h *Host) GetServices() {
	for i := range h.Services {
		fmt.Println(h.Services[i].ConfigMap)
	}

}
