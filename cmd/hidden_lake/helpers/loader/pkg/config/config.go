package config

import (
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/internal/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessagesCapacity: sett.GetMessagesCapacity(),
			FWorkSizeBits:     sett.GetWorkSizeBits(),
			FNetworkKey:       sett.GetNetworkKey(),
		},
	}
}
