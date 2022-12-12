package main

import (
	"github.com/number571/go-peer/cmd/hls/internal/config"
	hls_settings "github.com/number571/go-peer/cmd/hls/internal/settings"
	"github.com/number571/go-peer/modules/filesystem"
)

func initConfig() (config.IConfig, error) {
	if filesystem.OpenFile(hls_settings.CPathCFG).IsExist() {
		return config.LoadConfig(hls_settings.CPathCFG)
	}
	initCfg := &config.SConfig{
		FAddress: &config.SAddress{
			FTCP:  "localhost:9571",
			FHTTP: "localhost:9572",
		},
	}
	return config.NewConfig(hls_settings.CPathCFG, initCfg)
}
