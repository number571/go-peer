package config

import (
	"os"

	hlr_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/remoter/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	logger "github.com/number571/go-peer/internal/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FExecTimeoutMS: hlr_settings.CDefaultExecTimeout,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FIncoming: hlr_settings.CDefaultIncomingAddress,
				FPPROF:    "",
			},
			FConnection: hls_settings.CDefaultHTTPAddress,
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
