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
	FConnSettings conn.ISettings
	FAddress      string
	FMaxConnects  uint64
	FReadTimeout  time.Duration
	FWriteTimeout time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FConnSettings: pSett.FConnSettings,
		FAddress:      pSett.FAddress,
		FMaxConnects:  pSett.FMaxConnects,
		FReadTimeout:  pSett.FReadTimeout,
		FWriteTimeout: pSett.FWriteTimeout,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FConnSettings == nil {
		panic(`p.FConnSettings == nil`)
	}
	if p.FMaxConnects == 0 {
		panic(`p.FMaxConnects == 0`)
	}
	if p.FReadTimeout == 0 {
		panic(`p.FReadTimeout == 0`)
	}
	if p.FWriteTimeout == 0 {
		panic(`p.FWriteTimeout == 0`)
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
