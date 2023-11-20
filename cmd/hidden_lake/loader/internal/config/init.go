package config

import (
	"os"

	hll_settings "github.com/number571/go-peer/cmd/hidden_lake/loader/pkg/settings"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FProducers: []string{hll_settings.CDefaultProducerAddress},
			FConsumers: []string{hll_settings.CDefaultConsumerAddress},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
