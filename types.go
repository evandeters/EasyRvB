package main

import (
	"net"

	"github.com/lib/pq"
)

type Host struct {
	Hostname        string `gorm:"primaryKey"`
	Ip              net.IP
    NattedIp        net.IP
	VMTemplate      string
	Services        pq.StringArray `gorm:"type:text[]"`
}
