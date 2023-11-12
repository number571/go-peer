package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/filesystem"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if filesystem.OpenFile(cfgPath).IsExist() {
		cfg, err := LoadConfig(cfgPath)
		if err != nil {
			return nil, errors.WrapError(err, "load config by cfgPath")
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
				FMessagesCapacity:   hls_settings.CDefaultMessagesCapacity,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FTCP:  hls_settings.CDefaultTCPAddress,
				FHTTP: hls_settings.CDefaultHTTPAddress,
			},
			FServices: map[string]string{
				hlm_settings.CTitlePattern: hls_settings.CDefaultServiceHLMAddress,
			},
		}
	}
	cfg, err := BuildConfig(cfgPath, initCfg)
	if err != nil {
		return nil, errors.WrapError(err, "build config by cfgPath")
	}
	return cfg, nil
}
