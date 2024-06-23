package config

import (
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/config"
	"github.com/number571/go-peer/internal/language"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessagesCapacity: sett.GetMessagesCapacity(),
			FWorkSizeBits:     sett.GetWorkSizeBits(),
			FLanguage:         language.FromILanguage(sett.GetLanguage()),
		},
	}
}
