package queue

import (
	"time"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FMessageConstructSettings net_message.IConstructSettings
	FNetworkMask              uint32
	FConsumersCap             uint64
	FQueuePoolCap             [2]uint64
	FQueuePeriod              time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FMessageConstructSettings: pSett.FMessageConstructSettings,
		FNetworkMask:              pSett.FNetworkMask,
		FConsumersCap:             pSett.FConsumersCap,
		FQueuePoolCap:             pSett.FQueuePoolCap,
		FQueuePeriod:              pSett.FQueuePeriod,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FMessageConstructSettings == nil {
		panic(`p.FMessageConstructSettings == nil`)
	}
	if p.FQueuePeriod == 0 {
		panic(`p.FQueuePeriod == 0`)
	}
	if p.FQueuePoolCap[0] == 0 || p.FQueuePoolCap[1] == 0 {
		panic(`p.FQueuePoolCap[0] == 0 || p.FQueuePoolCap[1] == 0`)
	}
	if p.FConsumersCap == 0 {
		panic(`p.FConsumersCap == 0`)
	}
	// p.FNetworkMask can be = 0
	return p
}

func (p *sSettings) GetMessageConstructSettings() net_message.IConstructSettings {
	return p.FMessageConstructSettings
}

func (p *sSettings) GetNetworkMask() uint32 {
	return p.FNetworkMask
}

func (p *sSettings) GetQueuePoolCap() [2]uint64 {
	return p.FQueuePoolCap
}

func (p *sSettings) GetQueuePeriod() time.Duration {
	return p.FQueuePeriod
}

func (p *sSettings) GetConsumersCap() uint64 {
	return p.FConsumersCap
}
