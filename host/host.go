package host

import (
	"EasyRvB/service"
	"net"
)

type Host struct {
	Hostname        string
	Ip              net.IP
	OperatingSystem string
	Services        []*service.Service
}

func NewHost(hostname string, ip net.IP, operatingSystem string) *Host {
	return &Host{
		Hostname:        hostname,
		Ip:              ip,
		OperatingSystem: operatingSystem,
		Services:        []*service.Service{},
	}
}

func (h *Host) AddService(s service.Service) {
	h.Services = append(h.Services, &s)
}

func (h *Host) RemoveService(s service.Service) {
	for i, service := range h.Services {
		if service.Name == s.Name {
			h.Services = append(h.Services[:i], h.Services[i+1:]...)
			break
		}
	}
}
