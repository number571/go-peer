package queue

import (
	"time"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FMainCapacity       uint64
	FVoidCapacity       uint64
	FParallel           uint64
	FLimitVoidSizeBytes uint64
	FDuration           time.Duration
	FRandDuration       time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FMainCapacity:       pSett.FMainCapacity,
		FVoidCapacity:       pSett.FVoidCapacity,
		FParallel:           pSett.FParallel,
		FLimitVoidSizeBytes: pSett.FLimitVoidSizeBytes,
		FDuration:           pSett.FDuration,
		FRandDuration:       pSett.FRandDuration,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FMainCapacity == 0 {
		panic(`p.FMainCapacity == 0`)
	}
	if p.FVoidCapacity == 0 {
		panic(`p.FVoidCapacity == 0`)
	}
	if p.FParallel == 0 {
		panic(`p.FParallel == 0`)
	}
	if p.FDuration == 0 {
		panic(`p.FDuration == 0`)
	}
	// p.FRandDuration, p.FLimitVoidSizeBytes can be = 0
	return p
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
