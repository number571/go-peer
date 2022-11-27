package conn_keeper

import "time"

const (
	cTimeWait = time.Minute
)

type SSettings sSettings
type sSettings struct {
	FConnections func() []string
	FDuration    time.Duration
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FConnections: sett.FConnections,
		FDuration:    sett.FDuration,
	}).useDefaultValue()
}

func (s *sSettings) useDefaultValue() ISettings {
	if s.FDuration == 0 {
		s.FDuration = cTimeWait
	}
	return s
}

func (s *sSettings) GetConnections() []string {
	return s.FConnections()
}

func (s *sSettings) GetDuration() time.Duration {
	return s.FDuration
}
