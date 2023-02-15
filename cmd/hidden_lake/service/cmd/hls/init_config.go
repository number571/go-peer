package main

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/filesystem"
)

func initConfig() (config.IConfig, error) {
	if filesystem.OpenFile(pkg_settings.CPathCFG).IsExist() {
		return config.LoadConfig(pkg_settings.CPathCFG)
	}
	initCfg := &config.SConfig{
		FLogging: []string{config.CLogInfo, config.CLogWarn, config.CLogErro},
		FAddress: &config.SAddress{
			FTCP:  "localhost:9571",
			FHTTP: "localhost:9572",
		},
	}
	return config.NewConfig(pkg_settings.CPathCFG, initCfg)
}
