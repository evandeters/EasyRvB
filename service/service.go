package service

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
)

type ServiceConfig struct {
	Name       string
	Port       int16
	Protocol   string
	Limit      int
	Dependency map[string][]string

	Kubernetes KubernetesConfig `toml:"Kubernetes,omitempty"`
	Http       HTTPConfig       `toml:"Http,omitempty"`
}

type ServiceInstance struct {
	ID        uuid.UUID
	ConfigMap *ServiceConfig
}

func (g *ServiceConfig) ReadConfig(path string) error {
	if md, err := toml.DecodeFile(path, &g); err == nil {
		for _, undecoded := range md.Undecoded() {
			errMsg := fmt.Sprintf("[WARN] Undecoded configuration key: %v", undecoded)
			fmt.Println(errMsg)
		}
	} else {
		return err
	}
	return nil
}

func NewServiceInstance(config *ServiceConfig) *ServiceInstance {
	svc := ServiceInstance{
		ID:        uuid.New(),
		ConfigMap: config,
	}

	return &svc
}
