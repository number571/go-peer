package config

import (
	"os"

	hl_t_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/template/pkg/settings"
	logger "github.com/number571/go-peer/internal/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FValue: "TODO",
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FHTTP:  hl_t_settings.CDefaultHTTPAddress,
				FPPROF: "",
			},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
