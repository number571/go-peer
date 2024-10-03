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
	FMainPoolCapacity         uint64
	FRandPoolCapacity         uint64
	FQueuePeriod              time.Duration
	FRandQueuePeriod          time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FMessageConstructSettings: func() net_message.IConstructSettings {
			if pSett.FMessageConstructSettings == nil {
				return net_message.NewConstructSettings(&net_message.SConstructSettings{})
			}
			return pSett.FMessageConstructSettings
		}(),
		FNetworkMask:      pSett.FNetworkMask,
		FMainPoolCapacity: pSett.FMainPoolCapacity,
		FRandPoolCapacity: pSett.FRandPoolCapacity,
		FQueuePeriod:      pSett.FQueuePeriod,
		FRandQueuePeriod:  pSett.FRandQueuePeriod,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FMessageConstructSettings == nil {
		panic(`p.FMessageConstructSettings == nil`)
	}
	if p.FQueuePeriod == 0 {
		panic(`p.FQueuePeriod == 0`)
	}
	if p.FRandPoolCapacity == 0 {
		panic(`p.FRandPoolCapacity == 0`)
	}
	if p.FMainPoolCapacity == 0 {
		panic(`p.FMainPoolCapacity == 0`)
	}
	// p.FParallel, p.FNetworkMask, p.FWorkSizeBits, p.FRandQueuePeriod, p.FLimitVoidSizeBytes can be = 0
	return p
}

func (p *sSettings) GetMessageConstructSettings() net_message.IConstructSettings {
	return p.FMessageConstructSettings
}

func (p *sSettings) GetParallel() uint64 {
	return p.FMessageConstructSettings.GetParallel()
}

func (p *sSettings) GetRandMessageSizeBytes() uint64 {
	return p.FMessageConstructSettings.GetRandMessageSizeBytes()
}

func (p *sSettings) GetNetworkMask() uint32 {
	return p.FNetworkMask
}

func (p *sSettings) GetMainPoolCapacity() uint64 {
	return p.FMainPoolCapacity
}

func (p *sSettings) GetRandPoolCapacity() uint64 {
	return p.FRandPoolCapacity
}

func (p *sSettings) GetQueuePeriod() time.Duration {
	return p.FQueuePeriod
}

func (p *sSettings) GetRandQueuePeriod() time.Duration {
	return p.FRandQueuePeriod
}
