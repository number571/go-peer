package config

import (
	"os"

	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	hld_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/pkg/settings"
	logger "github.com/number571/go-peer/internal/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{},
			FLogging:  []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FHTTP:  hld_settings.CDefaultHTTPAddress,
				FPPROF: "",
			},
			FServices: map[string]string{
				hlm_settings.CServiceFullName: hlm_settings.CDefaultIncomingAddress,
				hlf_settings.CServiceFullName: hlf_settings.CDefaultIncomingAddress,
			},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
