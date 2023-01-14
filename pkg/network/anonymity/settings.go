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
	FTraffic      ITraffic
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FRetryEnqueue: sett.FRetryEnqueue,
		FNetworkMask:  sett.FNetworkMask,
		FTimeWait:     sett.FTimeWait,
		FTraffic:      sett.FTraffic,
	}).useDefaultValue()
}

func (s *sSettings) useDefaultValue() ISettings {
	if s.FNetworkMask == 0 {
		s.FNetworkMask = cMaskNetwork
	}
	if s.FTimeWait == 0 {
		s.FTimeWait = cTimeWait
	}
	if s.FTraffic == nil {
		s.FTraffic = NewTraffic(nil, nil)
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

func (s *sSettings) GetTraffic() ITraffic {
	return s.FTraffic
}
