package config

import (
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/internal/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessageSizeBytes:     sett.GetMessageSizeBytes(),
			FKeySizeBits:          sett.GetKeySizeBits(),
			FWorkSizeBits:         sett.GetWorkSizeBits(),
			FRandMessageSizeBytes: sett.GetRandMessageSizeBytes(),
			FNetworkKey:           sett.GetNetworkKey(),
		},
	}
}
