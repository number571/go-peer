package queue

import (
	"time"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FMainCapacity uint64
	FPoolCapacity uint64
	FParallel     uint64
	FDuration     time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FMainCapacity: pSett.FMainCapacity,
		FPoolCapacity: pSett.FPoolCapacity,
		FParallel:     pSett.FParallel,
		FDuration:     pSett.FDuration,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FMainCapacity == 0 {
		panic(`p.FMainCapacity == 0`)
	}
	if p.FPoolCapacity == 0 {
		panic(`p.FPoolCapacity == 0`)
	}
	if p.FParallel == 0 {
		panic(`p.FParallel == 0`)
	}
	if p.FDuration == 0 {
		panic(`p.FDuration == 0`)
	}
	return p
}

func (p *sSettings) GetParallel() uint64 {
	return p.FParallel
}

func (p *sSettings) GetMainCapacity() uint64 {
	return p.FMainCapacity
}

func (p *sSettings) GetPoolCapacity() uint64 {
	return p.FPoolCapacity
}

func (p *sSettings) GetDuration() time.Duration {
	return p.FDuration
}
