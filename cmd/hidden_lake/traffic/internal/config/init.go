package config

import (
	"github.com/number571/go-peer/internal/logger"
	"github.com/number571/go-peer/pkg/filesystem"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if filesystem.OpenFile(cfgPath).IsExist() {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FLogging:    []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress:    "localhost:9581",
			FConnection: "localhost:9571",
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
