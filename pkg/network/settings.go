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
	FAddress      string
	FMaxConnects  uint64
	FReadTimeout  time.Duration
	FWriteTimeout time.Duration
	FConnSettings conn.ISettings
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FAddress:      pSett.FAddress,
		FMaxConnects:  pSett.FMaxConnects,
		FReadTimeout:  pSett.FReadTimeout,
		FWriteTimeout: pSett.FWriteTimeout,
		FConnSettings: pSett.FConnSettings,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FMaxConnects == 0 {
		panic(`p.FMaxConnects == 0`)
	}
	if p.FReadTimeout == 0 {
		panic(`p.FReadTimeout == 0`)
	}
	if p.FWriteTimeout == 0 {
		panic(`p.FWriteTimeout == 0`)
	}
	if p.FConnSettings == nil {
		panic(`p.FConnSettings == nil`)
	}
	return p
}

func (p *sSettings) GetAddress() string {
	return p.FAddress
}

func (p *sSettings) GetMaxConnects() uint64 {
	return p.FMaxConnects
}

func (p *sSettings) GetConnSettings() conn.ISettings {
	return p.FConnSettings
}

func (p *sSettings) GetReadTimeout() time.Duration {
	return p.FReadTimeout
}

func (p *sSettings) GetWriteTimeout() time.Duration {
	return p.FWriteTimeout
}
