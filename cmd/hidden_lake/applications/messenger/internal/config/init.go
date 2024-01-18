package config

import (
	"os"

	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/crypto/random"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FMessagesCapacity: hlm_settings.CDefaultMessagesCapacity,
			},
			FLogging:   []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FLanguage:  hlm_settings.CDefaultLanguage,
			FShare:     hlm_settings.CDefaultShare,
			FPseudonym: random.NewStdPRNG().GetString(hlm_settings.CPseudonymSize),
			FAddress: &SAddress{
				FInterface: hlm_settings.CDefaultInterfaceAddress,
				FIncoming:  hlm_settings.CDefaultIncomingAddress,
			},
			FConnection: hls_settings.CDefaultHTTPAddress,
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
