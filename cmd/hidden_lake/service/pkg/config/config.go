package config

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/encoding"
)

func GetConfigSettings(pCfg config.IConfig, pClient client.IClient) SConfigSettings {
	sett := pCfg.GetSettings()
	msgLimit := pClient.GetMessageLimit()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessageSizeBytes:     sett.GetMessageSizeBytes(),
			FKeySizeBits:          sett.GetKeySizeBits(),
			FWorkSizeBits:         sett.GetWorkSizeBits(),
			FFetchTimeoutMS:       sett.GetFetchTimeoutMS(),
			FQueuePeriodMS:        sett.GetQueuePeriodMS(),
			FRandQueuePeriodMS:    sett.GetRandQueuePeriodMS(),
			FRandMessageSizeBytes: sett.GetRandMessageSizeBytes(),
			FTimestampWindowS:     sett.GetFetchTimeoutMS(),
			FNetworkKey:           sett.GetNetworkKey(),
		},
		// encoding.CSizeUint64 = payload64.Head()
		FLimitMessageSizeBytes: msgLimit - encoding.CSizeUint64,
	}
}
