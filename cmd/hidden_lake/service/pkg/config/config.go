package config

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func GetConfigSettings(pCfg config.IConfig, pNode anonymity.INode) SConfigSettings {
	sett := pCfg.GetSettings()
	msgLimit := pNode.GetMessageQueue().GetClient().GetMessageLimit()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessageSizeBytes:   sett.GetMessageSizeBytes(),
			FKeySizeBits:        sett.GetKeySizeBits(),
			FWorkSizeBits:       sett.GetWorkSizeBits(),
			FQueuePeriodMS:      sett.GetQueuePeriodMS(),
			FQueueRandPeriodMS:  sett.GetQueueRandPeriodMS(),
			FLimitVoidSizeBytes: sett.GetLimitVoidSizeBytes(),
			FNetworkKey:         sett.GetNetworkKey(),
			FF2FDisabled:        sett.GetF2FDisabled(),
			FQBTDisabled:        sett.GetQBTDisabled(),
		},
		// encoding.CSizeUint64 = payload64.Head()
		FLimitMessageSizeBytes: msgLimit - encoding.CSizeUint64,
	}
}
