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

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FConnections: pSett.FConnections,
		FDuration:    pSett.FDuration,
	}).useDefaultValue()
}

func (p *sSettings) useDefaultValue() ISettings {
	if p.FDuration == 0 {
		p.FDuration = cDuration
	}
	if p.FConnections == nil {
		p.FConnections = gConnections
	}
	return p
}

func (p *sSettings) GetConnections() []string {
	return p.FConnections()
}

func (p *sSettings) GetDuration() time.Duration {
	return p.FDuration
}
