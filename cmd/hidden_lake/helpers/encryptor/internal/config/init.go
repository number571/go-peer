package config

import (
	"os"

	hle_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/settings"
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
				FMessageSizeBytes:     hls_settings.CDefaultMessageSizeBytes,
				FWorkSizeBits:         hls_settings.CDefaultWorkSizeBits,
				FKeySizeBits:          hls_settings.CDefaultKeySizeBits,
				FRandMessageSizeBytes: hls_settings.CDefaultRandMessageSizeBytes,
				FNetworkKey:           hls_settings.CDefaultNetworkKey,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FHTTP:  hle_settings.CDefaultHTTPAddress,
				FPPROF: "",
			},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
