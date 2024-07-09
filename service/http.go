package service

type HTTPConfig struct {
	Host string `toml:"host"`
	RequireDB  bool `toml:"requiredb"`
}

func FillConfig(path string, values interface{}) error {
    return nil
}
