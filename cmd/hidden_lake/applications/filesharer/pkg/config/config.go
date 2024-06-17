package config

import (
	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/config"
	"github.com/number571/go-peer/internal/language"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FPageOffset: sett.GetPageOffset(),
			FRetryNum:   sett.GetRetryNum(),
			FLanguage:   language.FromILanguage(sett.GetLanguage()),
		},
	}
}
