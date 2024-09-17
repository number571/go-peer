package config

import (
	"fmt"
	"os"

	logger "github.com/number571/go-peer/internal/logger/std"

	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		cfg, err := LoadConfig(cfgPath)
		if err != nil {
			return nil, fmt.Errorf("load config: %w", err)
		}
		return cfg, nil
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FTimestampWindowS:     hls_settings.CDefaultTimestampWindowS,
				FMessageSizeBytes:     hls_settings.CDefaultMessageSizeBytes,
				FKeySizeBits:          hls_settings.CDefaultKeySizeBits,
				FWorkSizeBits:         hls_settings.CDefaultWorkSizeBits,
				FRandMessageSizeBytes: hls_settings.CDefaultRandMessageSizeBytes,
				FMessagesCapacity:     hlt_settings.CDefaultMessagesCapacity,
				FStorageEnabled:       hlt_settings.CDefaultStorageEnabled,
				FNetworkKey:           hls_settings.CDefaultNetworkKey,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FTCP:   hlt_settings.CDefaultTCPAddress,
				FHTTP:  hlt_settings.CDefaultHTTPAddress,
				FPPROF: "",
			},
			FConnections: []string{
				hlt_settings.CDefaultConnectionAddress,
			},
			FConsumers: []string{},
		}
	}
	cfg, err := BuildConfig(cfgPath, initCfg)
	if err != nil {
		return nil, fmt.Errorf("build config: %w", err)
	}
	return cfg, nil
}
