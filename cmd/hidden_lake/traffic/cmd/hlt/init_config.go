package main

import (
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/settings"
	"github.com/number571/go-peer/pkg/filesystem"
)

func initConfig() (config.IConfig, error) {
	if filesystem.OpenFile(hlt_settings.CPathCFG).IsExist() {
		return config.LoadConfig(hlt_settings.CPathCFG)
	}
	initCfg := &config.SConfig{
		FAddress:    "localhost:9573",
		FConnection: "localhost:9571",
	}
	return config.NewConfig(hlt_settings.CPathCFG, initCfg)
}
