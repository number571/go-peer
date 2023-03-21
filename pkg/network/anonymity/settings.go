package anonymity

import (
	"time"
)

var (
	_ ISettings = &sSettings{}
)

const (
	cServiceName = "DFT"
	cMaskNetwork = 0x1111111111111111
	cTimeWait    = time.Minute
)

type SSettings sSettings
type sSettings struct {
	FServiceName  string
	FRetryEnqueue uint64
	FNetworkMask  uint64
	FTimeWait     time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FServiceName:  pSett.FServiceName,
		FRetryEnqueue: pSett.FRetryEnqueue,
		FNetworkMask:  pSett.FNetworkMask,
		FTimeWait:     pSett.FTimeWait,
	}).useDefaultValue()
}

func (p *sSettings) useDefaultValue() ISettings {
	if p.FServiceName == "" {
		p.FServiceName = cServiceName
	}
	if p.FNetworkMask == 0 {
		p.FNetworkMask = cMaskNetwork
	}
	if p.FTimeWait == 0 {
		p.FTimeWait = cTimeWait
	}
	return p
}

func (p *sSettings) GetServiceName() string {
	return p.FServiceName
}

func (p *sSettings) GetTimeWait() time.Duration {
	return p.FTimeWait
}

func (p *sSettings) GetNetworkMask() uint64 {
	return p.FNetworkMask
}

func (p *sSettings) GetRetryEnqueue() uint64 {
	return p.FRetryEnqueue
}
