package host

import (
	"EasyRvB/service"
	"fmt"
	"net"
)

type Host struct {
	Hostname        string
	Ip              net.IP
	OperatingSystem string
	Services        []*service.ServiceInstance
}

func NewHost(hostname string, ip net.IP, operatingSystem string) *Host {
	return &Host{
		Hostname:        hostname,
		Ip:              ip,
		OperatingSystem: operatingSystem,
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
