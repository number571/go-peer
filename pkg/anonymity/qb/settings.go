package qb

import (
	"time"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FServiceName  string
	FFetchTimeout time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FServiceName:  pSett.FServiceName,
		FFetchTimeout: pSett.FFetchTimeout,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FFetchTimeout == 0 {
		panic(`p.FFetchTimeout == 0`)
	}
	return p
}

func (p *sSettings) GetServiceName() string {
	return p.FServiceName
}

func (p *sSettings) GetFetchTimeout() time.Duration {
	return p.FFetchTimeout
}
