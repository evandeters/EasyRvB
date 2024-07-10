package main

import (
	"net"
)

func NewHost(hostname string, ip, nattedIP net.IP, template string) *Host {
	return &Host{
		Hostname:        hostname,
		Ip:              ip,
        NattedIp:        nattedIP,
		VMTemplate:      template,
	}
}
