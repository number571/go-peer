package conn_keeper

import "time"

var (
	_ ISettings = &sSettings{}
)

const (
	cDuration = time.Minute
)

var (
	gConnections = func() []string { return nil }
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
		s.FDuration = cDuration
	}
	if s.FConnections == nil {
		s.FConnections = gConnections
	}
	return s
}

func (s *sSettings) GetConnections() []string {
	return s.FConnections()
}

func (s *sSettings) GetDuration() time.Duration {
	return s.FDuration
}
