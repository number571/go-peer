package config

import (
	"os"

	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/utils"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		cfg, err := LoadConfig(cfgPath)
		if err != nil {
			return nil, utils.MergeErrors(ErrLoadConfig, err)
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
				FQueueRandPeriodMS:  hls_settings.CDefaultQueueRandPeriod,
				FLimitVoidSizeBytes: hls_settings.CDefaultLimitVoidSize,
				FNetworkKey:         hls_settings.CDefaultNetworkKey,
				FF2FDisabled:        hls_settings.CDefaultF2FDisabled,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FTCP:  hls_settings.CDefaultTCPAddress,
				FHTTP: hls_settings.CDefaultHTTPAddress,
			},
			FServices: map[string]*SService{
				hlm_settings.CServiceFullName: {
					FHost: hlm_settings.CDefaultIncomingAddress,
				},
				hlf_settings.CServiceFullName: {
					FHost: hlf_settings.CDefaultIncomingAddress,
				},
			},
			FConnections: []string{},
			FFriends:     map[string]string{},
		}
	}
	cfg, err := BuildConfig(cfgPath, initCfg)
	if err != nil {
		return nil, utils.MergeErrors(ErrBuildConfig, err)
	}
	return cfg, nil
}
