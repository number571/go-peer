package config

import (
	"github.com/number571/go-peer/cmd/hms/utils"
)

var (
	_ IConfig = &sConfig{}
)

type sConfig struct {
	FAddress string `json:"address"`
}

const (
	// create local hms
	cAddress = "localhost:9572"
)

func NewConfig(filepath string) IConfig {
	var cfg = new(sConfig)

	if !utils.FileIsExist(filepath) {
		cfg = &sConfig{
			FAddress: cAddress,
		}
		err := utils.WriteFile(filepath, utils.Serialize(cfg))
		if err != nil {
			panic(err)
		}
	} else {
		err := utils.Deserialize(utils.ReadFile(filepath), cfg)
		if err != nil {
			panic(err)
		}
	}

	return cfg
}

func (cfg *sConfig) Address() string {
	return cfg.FAddress
}
