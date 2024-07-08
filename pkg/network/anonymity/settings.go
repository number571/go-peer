package anonymity

import (
	"time"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FServiceName  string
	FF2FDisabled  bool
	FNetworkMask  uint32
	FFetchTimeout time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FServiceName:  pSett.FServiceName,
		FF2FDisabled:  pSett.FF2FDisabled,
		FNetworkMask:  pSett.FNetworkMask,
		FFetchTimeout: pSett.FFetchTimeout,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FServiceName == "" {
		panic(`p.FServiceName == ""`)
	}
	if p.FNetworkMask == 0 {
		panic(`p.FNetworkMask == 0`)
	}
	if p.FFetchTimeout == 0 {
		panic(`p.FFetchTimeout == 0`)
	}
	return p
}

func (p *sSettings) GetF2FDisabled() bool {
	return p.FF2FDisabled
}

func (p *sSettings) GetServiceName() string {
	return p.FServiceName
}

func (p *sSettings) GetFetchTimeout() time.Duration {
	return p.FFetchTimeout
}

func (p *sSettings) GetNetworkMask() uint32 {
	return p.FNetworkMask
}
