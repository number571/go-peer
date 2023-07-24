package config

import (
	"github.com/number571/go-peer/internal/logger"
	"github.com/number571/go-peer/internal/settings"
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
			SConfigSettings: settings.SConfigSettings{
				FSettings: settings.SConfigSettingsBlock{
					FMessageSizeBytes: hls_settings.CDefaultMessageSize,
					FWorkSizeBits:     hls_settings.CDefaultWorkSize,
					FKeySizeBits:      hls_settings.CDefaultKeySize,
					FMessagesCapacity: hlm_settings.CDefaultCapMessages,
				},
			},
			FLogging:  []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FLanguage: "ENG",
			FAddress: &SAddress{
				FInterface: ":9591",
				FIncoming:  ":9592",
			},
			FConnection: &SConnection{
				FService: "service:9572",
				FTraffic: "traffic:9582",
			},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
