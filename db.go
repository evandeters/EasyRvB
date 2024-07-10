package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToSQLite() *gorm.DB {
    connString := "host=localhost user=easyrvb password=password dbname=easyrvb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
        panic(err)
	}
    return db
}

func AddHost(host Host) error {
    fmt.Println(host)
    result := db.Create(&host)
    if result.Error != nil {
        return result.Error
    }
    fmt.Println("Host added:", host.Hostname)
    return nil
}

func GetHosts() ([]Host, error) {
    var hosts []Host
    result := db.Find(&hosts)
    if result.Error != nil {
        return nil, result.Error
    }
    return hosts, nil
}

func AddService(host Host, service string) error {
    result := db.Model(&Host{}).Where("hostname = ?", host.Hostname).Update("services", gorm.Expr("array_append(services, ?)", service))
    if result.Error != nil {
        return result.Error
    }
    fmt.Println("Service added:", service)
    return nil
}

func GetServices(host Host) ([]string, error) {
    var hostRes Host
    result := db.Model(&Host{}).Where("hostname = ?", host.Hostname).Find(&hostRes)
    if result.Error != nil {
        return nil, result.Error
    }
    return hostRes.Services, nil
}

func GetHostsWithService(service string) ([]Host, error) {
    var hosts []Host
    result := db.Model(&Host{}).Where("? = ANY (services)", service).Find(&hosts)
    if result.Error != nil {
        return nil, result.Error
    }
    return hosts, nil
}
