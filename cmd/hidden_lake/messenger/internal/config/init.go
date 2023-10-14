package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/filesystem"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if filesystem.OpenFile(cfgPath).IsExist() {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FMessagesCapacity: hlm_settings.CDefaultCapMessages,
				FWorkSizeBits:     hlm_settings.CDefaultWorkSizeBits,
			},
			FLogging:  []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FLanguage: hlm_settings.CDefaultLanguage,
			FAddress: &SAddress{
				FInterface: hlm_settings.CDefaultInterfaceAddress,
				FIncoming:  hlm_settings.CDefaultIncomingAddress,
			},
			FConnection: hls_settings.CDefaultHTTPAddress,
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
