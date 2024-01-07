package config

import (
	"fmt"
	"os"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	logger "github.com/number571/go-peer/internal/logger/std"
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
				FKeySizeBits:        hls_settings.CDefaultKeySize,
				FQueuePeriodMS:      hls_settings.CDefaultQueuePeriod,
				FLimitVoidSizeBytes: hls_settings.CDefaultLimitVoidSize,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FTCP:  hls_settings.CDefaultTCPAddress,
				FHTTP: hls_settings.CDefaultHTTPAddress,
			},
			FF2FDisabled: hls_settings.CDefaultF2FDisabled,
			FServices: map[string]*SService{
				hlm_settings.CTitlePattern: {FHost: hls_settings.CDefaultServiceHLMAddress},
			},
		}
	}
	cfg, err := BuildConfig(cfgPath, initCfg)
	if err != nil {
		return nil, fmt.Errorf("build config: %w", err)
	}
	return cfg, nil
}
