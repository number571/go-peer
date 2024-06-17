package config

import (
	"github.com/number571/go-peer/cmd/hidden_lake/template/internal/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FValue: sett.GetValue(),
		},
	}
}
