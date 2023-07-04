package network

import (
	"time"

	"github.com/number571/go-peer/pkg/network/conn"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FAddress       string
	FCapacity      uint64
	FMaxConnects   uint64
	FActionTimeout time.Duration
	FConnSettings  conn.ISettings
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FAddress:       pSett.FAddress,
		FCapacity:      pSett.FCapacity,
		FMaxConnects:   pSett.FMaxConnects,
		FActionTimeout: pSett.FActionTimeout,
		FConnSettings:  pSett.FConnSettings,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FCapacity == 0 {
		panic(`p.FCapacity == 0`)
	}
	if p.FMaxConnects == 0 {
		panic(`p.FMaxConnects == 0`)
	}
	if p.FActionTimeout == 0 {
		panic(`p.FActionTimeout == 0`)
	}
	if p.FConnSettings == nil {
		panic(`p.FConnSettings == nil`)
	}
	return p
}

func (p *sSettings) GetAddress() string {
	return p.FAddress
}

func (p *sSettings) GetCapacity() uint64 {
	return p.FCapacity
}

func (p *sSettings) GetMaxConnects() uint64 {
	return p.FMaxConnects
}

func (p *sSettings) GetConnSettings() conn.ISettings {
	return p.FConnSettings
}

func (p *sSettings) GetActionTimeout() time.Duration {
	return p.FActionTimeout
}
