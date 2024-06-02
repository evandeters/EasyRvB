package service

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type KubernetesConfig struct {
	GenericService
	Distribution string
	NodeType     string
	Dependency   map[string]any
}

func (k *KubernetesConfig) ReadConfig(path string) error {
	if md, err := toml.DecodeFile(path, &k); err != nil {
		return err
	} else {
		for _, undecoded := range md.Undecoded() {
			errMsg := fmt.Sprintf("[WARN] Undecoded configuration key: %v", undecoded)
			fmt.Println(errMsg)
		}
	}
	return nil
}
