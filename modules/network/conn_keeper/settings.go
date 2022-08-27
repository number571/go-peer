package conn_keeper

import "time"

const (
	cTimeWait = time.Minute
)

type SSettings struct {
	FConnections []string
	FDuration    time.Duration
}

func NewSettings(sett *SSettings) ISettings {
	return (&SSettings{
		FConnections: sett.FConnections,
		FDuration:    sett.FDuration,
	}).useDefaultValue()
}

func (s *SSettings) useDefaultValue() ISettings {
	if s.FDuration == 0 {
		s.FDuration = cTimeWait
	}
	return s
}

func (s *SSettings) GetConnections() []string {
	return s.FConnections
}

func (s *SSettings) GetDuration() time.Duration {
	return s.FDuration
}
