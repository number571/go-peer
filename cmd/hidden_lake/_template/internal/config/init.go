package config

import (
	"os"

	hl_t_settings "github.com/number571/go-peer/cmd/hidden_lake/_template/pkg/settings"
	logger "github.com/number571/go-peer/internal/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FHTTP: hl_t_settings.CDefaultHTTPAddress,
			},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
