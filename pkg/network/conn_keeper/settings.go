package conn_keeper // nolint: revive

import "time"

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FConnections func() []string
	FDuration    time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FConnections: pSett.FConnections,
		FDuration:    pSett.FDuration,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FDuration == 0 {
		panic(`p.FDuration = cDuration`)
	}
	if p.FConnections == nil {
		panic(`p.FConnections == nil`)
	}
	return p
}

func (p *sSettings) GetConnections() []string {
	return p.FConnections()
}

func (p *sSettings) GetDuration() time.Duration {
	return p.FDuration
}
