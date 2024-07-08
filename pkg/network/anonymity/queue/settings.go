package queue

import (
	"time"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FNetworkMask          uint32
	FWorkSizeBits         uint64
	FMainPoolCapacity     uint64
	FRandPoolCapacity     uint64
	FParallel             uint64
	FRandMessageSizeBytes uint64
	FQueuePeriod          time.Duration
	FRandQueuePeriod      time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FNetworkMask:          pSett.FNetworkMask,
		FWorkSizeBits:         pSett.FWorkSizeBits,
		FMainPoolCapacity:     pSett.FMainPoolCapacity,
		FRandPoolCapacity:     pSett.FRandPoolCapacity,
		FParallel:             pSett.FParallel,
		FRandMessageSizeBytes: pSett.FRandMessageSizeBytes,
		FQueuePeriod:          pSett.FQueuePeriod,
		FRandQueuePeriod:      pSett.FRandQueuePeriod,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FQueuePeriod != 0 && p.FRandPoolCapacity == 0 {
		panic(`p.FQueuePeriod != 0 && p.FRandPoolCapacity == 0`)
	}
	if p.FMainPoolCapacity == 0 {
		panic(`p.FMainPoolCapacity == 0`)
	}
	// p.FParallel, p.FNetworkMask, p.FWorkSizeBits, p.FRandQueuePeriod, p.FLimitVoidSizeBytes can be = 0
	return p
}

func (p *sSettings) GetNetworkMask() uint32 {
	return p.FNetworkMask
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *sSettings) GetParallel() uint64 {
	return p.FParallel
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

func (p *sSettings) GetRandMessageSizeBytes() uint64 {
	return p.FRandMessageSizeBytes
}
