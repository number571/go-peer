package network

import "time"

const (
	cNetworkKey  = "network-key"
	cCapacity    = (1 << 10)
	cMessageSize = (1 << 20)
	cMaxConnects = (1 << 6)
	cTimeWait    = time.Minute
)

type SSettings struct {
	FNetworkKey  string
	FCapacity    uint64
	FMessageSize uint64
	FMaxConnects uint64
	FTimeWait    time.Duration
}

func NewSettings(sett *SSettings) ISettings {
	return (&SSettings{
		FNetworkKey:  sett.FNetworkKey,
		FCapacity:    sett.FCapacity,
		FMessageSize: sett.FMessageSize,
		FMaxConnects: sett.FMaxConnects,
		FTimeWait:    sett.FTimeWait,
	}).useDefaultValues()
}

func (s *SSettings) useDefaultValues() ISettings {
	if s.FNetworkKey == "" {
		s.FNetworkKey = cNetworkKey
	}
	if s.FCapacity == 0 {
		s.FCapacity = cCapacity
	}
	if s.FMessageSize == 0 {
		s.FMessageSize = cMessageSize
	}
	if s.FMaxConnects == 0 {
		s.FMaxConnects = cMaxConnects
	}
	if s.FTimeWait == 0 {
		s.FTimeWait = cTimeWait
	}
	return s
}

func (s *SSettings) GetNetworkKey() string {
	return s.FNetworkKey
}

func (s *SSettings) GetCapacity() uint64 {
	return s.FCapacity
}

func (s *SSettings) GetMessageSize() uint64 {
	return s.FMessageSize
}

func (s *SSettings) GetMaxConnects() uint64 {
	return s.FMaxConnects
}

func (s *SSettings) GetTimeWait() time.Duration {
	return s.FTimeWait
}
