package anonymity

import (
	"time"
)

var (
	_ ISettings = &sSettings{}
)

const (
	cMaskNetwork = 0x676F2D7065657201
	cTimeWait    = time.Minute
)

type SSettings sSettings
type sSettings struct {
	FRetryEnqueue uint64
	FNetworkMask  uint64
	FTimeWait     time.Duration
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FRetryEnqueue: sett.FRetryEnqueue,
		FNetworkMask:  sett.FNetworkMask,
		FTimeWait:     sett.FTimeWait,
	}).useDefaultValue()
}

func (s *sSettings) useDefaultValue() ISettings {
	if s.FNetworkMask == 0 {
		s.FNetworkMask = cMaskNetwork
	}
	if s.FTimeWait == 0 {
		s.FTimeWait = cTimeWait
	}
	return s
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
