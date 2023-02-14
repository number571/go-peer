package config

import (
	"fmt"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/filesystem"
)

var (
	_ IConfig = &SConfig{}
)

type SConfig struct {
	FNetwork    string `json:"network,omitempty"`
	FAddress    string `json:"address"`
	FConnection string `json:"connection"`
}

func NewConfig(filepath string, cfg *SConfig) (IConfig, error) {
	configFile := filesystem.OpenFile(filepath)

	if configFile.IsExist() {
		return nil, fmt.Errorf("config file '%s' already exist", filepath)
	}

	if err := configFile.Write(encoding.Serialize(cfg)); err != nil {
		return nil, err
	}

	return cfg, nil
}

func LoadConfig(filepath string) (IConfig, error) {
	configFile := filesystem.OpenFile(filepath)

	if !configFile.IsExist() {
		return nil, fmt.Errorf("config file '%s' does not exist", filepath)
	}

	bytes, err := configFile.Read()
	if err != nil {
		return nil, err
	}

	cfg := new(SConfig)
	if err := encoding.Deserialize(bytes, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *SConfig) Network() string {
	return cfg.FNetwork
}

func (cfg *SConfig) Address() string {
	return cfg.FAddress
}

func (cfg *SConfig) Connection() string {
	return cfg.FConnection
}
