package main

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/internal/settings"
	"github.com/number571/go-peer/pkg/filesystem"
)

func initConfig() (config.IConfig, error) {
	if filesystem.OpenFile(hls_settings.CPathCFG).IsExist() {
		return config.LoadConfig(hls_settings.CPathCFG)
	}
	initCfg := &config.SConfig{
		FLogging: []string{config.CLogInfo, config.CLogWarn, config.CLogErro},
		FAddress: &config.SAddress{
			FTCP:  "localhost:9571",
			FHTTP: "localhost:9572",
		},
	}
	return config.NewConfig(hls_settings.CPathCFG, initCfg)
}
