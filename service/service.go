package service

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Service interface {
    SetConfig(resPath string) error
}

type ServiceConfig struct {
	Name       string `toml:"Name"`
    User       string `toml:"User"`
    Type       string `toml:"Type"`
	Port       int16  `toml:"Port"`
	Protocol   string `toml:"Protocol"`
	Limit      int    `toml:"Limit"`
    SupportedOS []string `toml:"SupportedOS,omitempty"`
	Dependency []string `toml:"Dependency,omitempty"`

	Kubernetes KubernetesConfig `toml:"Kubernetes,omitempty"`
	Http       *HTTPConfig       `toml:"Http,omitempty"`
    Database   DatabaseConfig      `toml:"Database,omitempty"`
    Dns        DnsConfig           `toml:"Dns,omitempty"`
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
