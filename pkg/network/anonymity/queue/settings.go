package queue

import (
	"time"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FNetworkMask        uint32
	FQBTDisabled        bool
	FWorkSizeBits       uint64
	FMainCapacity       uint64
	FVoidCapacity       uint64
	FParallel           uint64
	FLimitVoidSizeBytes uint64
	FDuration           time.Duration
	FRandDuration       time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FNetworkMask:        pSett.FNetworkMask,
		FQBTDisabled:        pSett.FQBTDisabled,
		FWorkSizeBits:       pSett.FWorkSizeBits,
		FMainCapacity:       pSett.FMainCapacity,
		FVoidCapacity:       pSett.FVoidCapacity,
		FParallel:           pSett.FParallel,
		FLimitVoidSizeBytes: pSett.FLimitVoidSizeBytes,
		FDuration:           pSett.FDuration,
		FRandDuration:       pSett.FRandDuration,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if !p.FQBTDisabled && p.FVoidCapacity == 0 {
		panic(`!p.FQBTDisabled && p.FVoidCapacity == 0`)
	}
	if p.FMainCapacity == 0 {
		panic(`p.FMainCapacity == 0`)
	}
	if p.FDuration == 0 {
		panic(`p.FDuration == 0`)
	}
	// p.FParallel, p.FNetworkMask, p.FWorkSizeBits, p.FRandDuration, p.FLimitVoidSizeBytes can be = 0
	return p
}

func (p *sSettings) GetNetworkMask() uint32 {
	return p.FNetworkMask
}

func (p *sSettings) GetQBTDisabled() bool {
	return p.FQBTDisabled
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *sSettings) GetParallel() uint64 {
	return p.FParallel
}

func (p *sSettings) GetMainCapacity() uint64 {
	return p.FMainCapacity
}

func (p *sSettings) GetVoidCapacity() uint64 {
	return p.FVoidCapacity
}

func (p *sSettings) GetDuration() time.Duration {
	return p.FDuration
}

func (p *sSettings) GetRandDuration() time.Duration {
	return p.FRandDuration
}

func (p *sSettings) GetLimitVoidSizeBytes() uint64 {
	return p.FLimitVoidSizeBytes
}
