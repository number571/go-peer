package config

import (
	"github.com/number571/go-peer/internal/logger"
	"github.com/number571/go-peer/internal/settings"
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
			SConfigSettings: settings.SConfigSettings{
				FSettings: settings.SConfigSettingsBlock{
					FMessageSize: hls_settings.CDefaultMessageSize,
					FWorkSize:    hls_settings.CDefaultWorkSize,
					FKeySize:     hls_settings.CDefaultKeySize,
					FQueuePeriod: hls_settings.CDefaultQueuePeriod,
				},
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FTCP:  ":9571",
				FHTTP: ":9572",
			},
			FServices: map[string]string{
				hlm_settings.CTitlePattern: "messenger:9592",
			},
		}
	}
	cfg, err := BuildConfig(cfgPath, initCfg)
	if err != nil {
		return nil, errors.WrapError(err, "build config by cfgPath")
	}
	return cfg, nil
}
