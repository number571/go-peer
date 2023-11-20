package config

import (
	"fmt"
	"os"

	logger "github.com/number571/go-peer/internal/logger/std"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
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
				FMessageSizeBytes:   hls_settings.CDefaultMessageSize,
				FWorkSizeBits:       hls_settings.CDefaultWorkSize,
				FQueuePeriodMS:      hls_settings.CDefaultQueuePeriod,
				FLimitVoidSizeBytes: hls_settings.CDefaultLimitVoidSize,
				FMessagesCapacity:   hlt_settings.CDefaultMessagesCapacity,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FTCP:  hlt_settings.CDefaultTCPAddress,
				FHTTP: hlt_settings.CDefaultHTTPAddress,
			},
			FConnections: []string{
				hlt_settings.CDefaultConnectionAddress,
			},
		}
	}
	cfg, err := BuildConfig(cfgPath, initCfg)
	if err != nil {
		return nil, fmt.Errorf("build config: %w", err)
	}
	return cfg, nil
}
