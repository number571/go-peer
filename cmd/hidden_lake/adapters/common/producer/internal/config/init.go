package config

import (
	"os"

	hla_settings "github.com/number571/go-peer/cmd/hidden_lake/adapters/common/producer/pkg/settings"
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
				FWorkSizeBits: hls_settings.CDefaultWorkSize,
				FNetworkKey:   hls_settings.CDefaultNetworkKey,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: hla_settings.CDefaultHTTPAddress,
			FConnection: &SConnection{
				FSrvHost: hla_settings.CDefaultSrvAddress,
			},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
