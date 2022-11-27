package conn

import "time"

const (
	cNetworkKey  = "network-key"
	cMessageSize = (1 << 20)
	cTimeWait    = time.Minute
)

type SSettings sSettings
type sSettings struct {
	FNetworkKey  string
	FMessageSize uint64
	FTimeWait    time.Duration
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FNetworkKey:  sett.FNetworkKey,
		FMessageSize: sett.FMessageSize,
		FTimeWait:    sett.FTimeWait,
	}).useDefaultValues()
}

func (s *sSettings) useDefaultValues() ISettings {
	if s.FNetworkKey == "" {
		s.FNetworkKey = cNetworkKey
	}
	if s.FMessageSize == 0 {
		s.FMessageSize = cMessageSize
	}
	if s.FTimeWait == 0 {
		s.FTimeWait = cTimeWait
	}
	return s
}

func (s *sSettings) GetNetworkKey() string {
	return s.FNetworkKey
}

func (s *sSettings) GetMessageSize() uint64 {
	return s.FMessageSize
}

func (s *sSettings) GetTimeWait() time.Duration {
	return s.FTimeWait
}
