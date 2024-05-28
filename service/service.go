package service

import (
    "net"
    "encoding/json"
    "fmt"
)

type Service struct {
    Name string
    Port int
    Protocol string
    Dependencies []Dependency
    Config ServiceConfig
}

type ServiceConfig interface {
    ReadConfig() error
}

type Dependency struct {
    Name string
    Ip net.IP
    Port int16
}

func NewService(name string, port int, protocol string, dependencies []Dependency, config ServiceConfig) *Service {
    return &Service{
        Name: name,
        Port: port,
        Protocol: protocol,
        Dependencies: dependencies,
        Config: nil,
    }
}

func (s *Service) AddDependency(d Dependency) {
    s.Dependencies = append(s.Dependencies, d)
}

func (s *Service) RemoveDependency(d Dependency) {
    for i, dependency := range s.Dependencies {
        if dependency.Name == d.Name {
            s.Dependencies = append(s.Dependencies[:i], s.Dependencies[i+1:]...)
            break
        }
    }
}

func NewServiceConfig(serviceType string, data []byte) (ServiceConfig, error) {
    var config ServiceConfig
    switch serviceType {
    case "HTTP":
        config = &HTTPConfig{}
    default:
        return nil, fmt.Errorf("Service type %s not supported", serviceType)
    }

    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return config, nil
}
