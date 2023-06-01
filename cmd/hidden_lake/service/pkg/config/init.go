package config

import (
	"github.com/number571/go-peer/internal/logger"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/filesystem"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if filesystem.OpenFile(cfgPath).IsExist() {
		cfg, err := LoadConfig(cfgPath)
		if err != nil {
			return nil, errors.WrapError(err, "load config by cfgPath")
		}
		return cfg, nil
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FTCP:  "localhost:9571",
				FHTTP: "localhost:9572",
			},
		}
	}
	cfg, err := BuildConfig(cfgPath, initCfg)
	if err != nil {
		return nil, errors.WrapError(err, "build config by cfgPath")
	}
	return cfg, nil
}
