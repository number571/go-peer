package main

import (
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/pkg/filesystem"
)

func initConfig() (config.IConfig, error) {
	if filesystem.OpenFile(hlm_settings.CPathCFG).IsExist() {
		return config.LoadConfig(hlm_settings.CPathCFG)
	}
	initCfg := &config.SConfig{
		FAddress: &config.SAddress{
			FInterface: "localhost:8080",
			FIncoming:  "localhost:8081",
		},
		FConnection: &config.SConnection{
			FService: "localhost:9572",
		},
	}
	return config.NewConfig(hlm_settings.CPathCFG, initCfg)
}
