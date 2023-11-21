package config

import (
	"os"

	hll_settings "github.com/number571/go-peer/cmd/hidden_lake/loader/pkg/settings"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	logger "github.com/number571/go-peer/internal/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FMessagesCapacity: hlt_settings.CDefaultMessagesCapacity,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FHTTP: hll_settings.CDefaultHTTPAddress,
			},
			FProducers: []string{hll_settings.CDefaultProducerAddress},
			FConsumers: []string{hll_settings.CDefaultConsumerAddress},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
