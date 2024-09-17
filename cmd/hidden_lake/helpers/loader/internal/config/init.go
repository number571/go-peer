package config

import (
	"os"

	hll_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/settings"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/settings"
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
				FMessagesCapacity: hlt_settings.CDefaultMessagesCapacity,
				FWorkSizeBits:     hls_settings.CDefaultWorkSizeBits,
				FNetworkKey:       hls_settings.CDefaultNetworkKey,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FHTTP:  hll_settings.CDefaultHTTPAddress,
				FPPROF: "",
			},
			FProducers: []string{},
			FConsumers: []string{},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
