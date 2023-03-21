package queue

import (
	"time"
)

var (
	_ ISettings = &sSettings{}
)

const (
	cCapacity     = (1 << 5)
	cPullCapacity = (1 << 5)
	cDuration     = time.Second
)

type SSettings sSettings
type sSettings struct {
	FCapacity     uint64
	FPullCapacity uint64
	FDuration     time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FCapacity:     pSett.FCapacity,
		FPullCapacity: pSett.FPullCapacity,
		FDuration:     pSett.FDuration,
	}).useDefaultValues()
}

func (p *sSettings) useDefaultValues() ISettings {
	if p.FCapacity == 0 {
		p.FCapacity = cCapacity
	}
	if p.FPullCapacity == 0 {
		p.FPullCapacity = cPullCapacity
	}
	if p.FDuration == 0 {
		p.FDuration = cDuration
	}
	return p
}

func (p *sSettings) GetCapacity() uint64 {
	return p.FCapacity
}

func (p *sSettings) GetPullCapacity() uint64 {
	return p.FPullCapacity
}

func (p *sSettings) GetDuration() time.Duration {
	return p.FDuration
}
