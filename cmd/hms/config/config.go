package config

import (
	"github.com/number571/go-peer/cmd/hms/utils"
)

var (
	_ IConfig = &sConfig{}
)

type sConfig struct {
	FAddress   string `json:"address"`
	FCleanCron string `json:"clean_cron"`
}

const (
	// create local hms
	cDefaultAddress = "localhost:8080"

	// cron of clean database
	cDefaultCleanCron = "0 0 * * *"
)

func NewConfig(filepath string) IConfig {
	var cfg = new(sConfig)

	if !utils.FileIsExist(filepath) {
		cfg = &sConfig{
			FAddress:   cDefaultAddress,
			FCleanCron: cDefaultCleanCron,
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

func (cfg *sConfig) CleanCron() string {
	return cfg.FCleanCron
}
