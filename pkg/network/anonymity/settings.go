package anonymity

import (
	"time"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FServiceName   string
	FRetryEnqueue  uint64
	FNetworkMask   uint64
	FFetchTimeWait time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FServiceName:   pSett.FServiceName,
		FRetryEnqueue:  pSett.FRetryEnqueue,
		FNetworkMask:   pSett.FNetworkMask,
		FFetchTimeWait: pSett.FFetchTimeWait,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FServiceName == "" {
		panic(`p.FServiceName == ""`)
	}
	if p.FNetworkMask == 0 {
		panic(`p.FNetworkMask == 0`)
	}
	if p.FFetchTimeWait == 0 {
		panic(`p.FTimeWait == 0`)
	}
	return p
}

func (p *sSettings) GetServiceName() string {
	return p.FServiceName
}

func (p *sSettings) GetFetchTimeWait() time.Duration {
	return p.FFetchTimeWait
}

func (p *sSettings) GetNetworkMask() uint64 {
	return p.FNetworkMask
}

func (p *sSettings) GetRetryEnqueue() uint64 {
	return p.FRetryEnqueue
}
