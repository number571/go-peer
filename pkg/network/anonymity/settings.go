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

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FServiceName:  sett.FServiceName,
		FRetryEnqueue: sett.FRetryEnqueue,
		FNetworkMask:  sett.FNetworkMask,
		FTimeWait:     sett.FTimeWait,
	}).useDefaultValue()
}

func (s *sSettings) useDefaultValue() ISettings {
	if s.FServiceName == "" {
		s.FServiceName = cServiceName
	}
	if s.FNetworkMask == 0 {
		s.FNetworkMask = cMaskNetwork
	}
	if s.FTimeWait == 0 {
		s.FTimeWait = cTimeWait
	}
	return s
}

func (s *sSettings) GetServiceName() string {
	return s.FServiceName
}

func (s *sSettings) GetTimeWait() time.Duration {
	return s.FTimeWait
}

func (s *sSettings) GetNetworkMask() uint64 {
	return s.FNetworkMask
}

func (s *sSettings) GetRetryEnqueue() uint64 {
	return s.FRetryEnqueue
}
