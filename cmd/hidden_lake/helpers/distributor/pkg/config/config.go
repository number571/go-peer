package config

import (
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/internal/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	_ = pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{},
	}
}
