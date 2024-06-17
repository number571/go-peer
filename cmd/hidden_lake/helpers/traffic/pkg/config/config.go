package config

import (
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessageSizeBytes:   sett.GetMessageSizeBytes(),
			FKeySizeBits:        sett.GetKeySizeBits(),
			FWorkSizeBits:       sett.GetWorkSizeBits(),
			FMessagesCapacity:   sett.GetMessagesCapacity(),
			FLimitVoidSizeBytes: sett.GetLimitVoidSizeBytes(),
			FNetworkKey:         sett.GetNetworkKey(),
			FStorageEnabled:     sett.GetStorageEnabled(),
		},
	}
}
