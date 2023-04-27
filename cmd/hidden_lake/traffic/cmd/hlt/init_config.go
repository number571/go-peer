package main

import (
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/logger"
	"github.com/number571/go-peer/pkg/filesystem"
)

func initConfig() (config.IConfig, error) {
	if filesystem.OpenFile(hlt_settings.CPathCFG).IsExist() {
		return config.LoadConfig(hlt_settings.CPathCFG)
	}
	initCfg := &config.SConfig{
		FLogging:    []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
		FAddress:    "localhost:9581",
		FConnection: "localhost:9571",
	}
	return config.BuildConfig(hlt_settings.CPathCFG, initCfg)
}
