package config

import (
	"github.com/number571/go-peer/cmd/hidden_lake/applications/remoter/internal/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FExecTimeoutMS: sett.GetExecTimeoutMS(),
		},
	}
}
