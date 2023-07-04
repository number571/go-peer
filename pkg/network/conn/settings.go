package conn

import "time"

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FNetworkKey    string
	FMessageSize   uint64
	FLimitVoidSize uint64
	FReadDeadline  time.Duration
	FWriteDeadline time.Duration
	FFetchTimeWait time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FNetworkKey:    pSett.FNetworkKey,
		FMessageSize:   pSett.FMessageSize,
		FLimitVoidSize: pSett.FLimitVoidSize,
		FReadDeadline:  pSett.FReadDeadline,
		FWriteDeadline: pSett.FWriteDeadline,
		FFetchTimeWait: pSett.FFetchTimeWait,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FMessageSize == 0 {
		panic(`p.FMessageSize == 0`)
	}
	if p.FLimitVoidSize == 0 {
		panic(`p.FMaxVoidSize == 0`)
	}
	if p.FReadDeadline == 0 {
		panic(`p.FReadDeadline == 0`)
	}
	if p.FWriteDeadline == 0 {
		panic(`p.FWriteDeadline == 0`)
	}
	if p.FFetchTimeWait == 0 {
		panic(`p.FFetchTimeWait == 0`)
	}
	return p
}

func (p *sSettings) GetNetworkKey() string {
	return p.FNetworkKey
}

func (p *sSettings) GetMessageSize() uint64 {
	return p.FMessageSize
}

func (p *sSettings) GetLimitVoidSize() uint64 {
	return p.FLimitVoidSize
}

func (p *sSettings) GetReadDeadline() time.Duration {
	return p.FReadDeadline
}

func (p *sSettings) GetWriteDeadline() time.Duration {
	return p.FWriteDeadline
}

func (p *sSettings) GetFetchTimeWait() time.Duration {
	return p.FFetchTimeWait
}
