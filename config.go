package main

import (
    "log"
    "os"

    "github.com/BurntSushi/toml"
    "github.com/pkg/errors"
)

var (
    configErrors = []string{}
)

type Config struct {
    VCenterServer              string
    VCenterUsername            string
    VCenterPassword            string
    Datacenter                 string
    Datastore                  string
    Cluster                    string
    VMFolder                   string
    ResourcePool               string
    LogPath                    string
    PortGroup                  string
    RouterTemplate             string
    VMTemplates                []string
}

/*
Load config settings into given config object
*/
func ReadConfig(conf *Config, configPath string) {
    fileContent, err := os.ReadFile(configPath)
    if err != nil {
        log.Fatalln("Configuration file ("+configPath+") not found:", err)
    }
    if md, err := toml.Decode(string(fileContent), &conf); err != nil {
        log.Fatalln(err)
    } else {
        for _, undecoded := range md.Undecoded() {
            errMsg := "[WARN] Undecoded configuration key \"" + undecoded.String() + "\" will not be used."
            configErrors = append(configErrors, errMsg)
            log.Println(errMsg)
        }
    }
}

/*
Check for config errors and set defaults
*/
func CheckConfig(conf *Config) error {
    if conf.VCenterServer == "" {
        return errors.New("illegal config: vCenterURL must be defined")
    }
    if conf.VCenterUsername == "" {
        return errors.New("illegal config: vCenterUsername must be defined")
    }
    if conf.VCenterPassword == "" {
        return errors.New("illegal config: vCenterPassword must be defined")
    }
    if conf.Datacenter == "" {
        return errors.New("illegal config: Datacenter must be defined")
    }
    if conf.ResourcePool == "" {
        return errors.New("illegal config: PresetTemplateResourcePool must be defined")
    }
    if conf.PortGroup == "" {
        return errors.New("illegal config: MainDistributedSwitch must be defined")
    }

    return nil
}
