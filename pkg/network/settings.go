package network

import "github.com/number571/go-peer/pkg/network/conn"

var (
	_ ISettings = &sSettings{}
)

const (
	cAddress     = ""
	cCapacity    = (1 << 10)
	cMaxConnects = (1 << 6)
)

type SSettings sSettings
type sSettings struct {
	FAddress      string
	FCapacity     uint64
	FMaxConnects  uint64
	FConnSettings conn.ISettings
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FAddress:      pSett.FAddress,
		FCapacity:     pSett.FCapacity,
		FMaxConnects:  pSett.FMaxConnects,
		FConnSettings: pSett.FConnSettings,
	}).useDefaultValues()
}

func (p *sSettings) useDefaultValues() ISettings {
	if p.FAddress == "" {
		p.FAddress = cAddress
	}
	if p.FCapacity == 0 {
		p.FCapacity = cCapacity
	}
	if p.FMaxConnects == 0 {
		p.FMaxConnects = cMaxConnects
	}
	if p.FConnSettings == nil {
		p.FConnSettings = conn.NewSettings(&conn.SSettings{})
	}
	return p
}

func (p *sSettings) GetAddress() string {
	return p.FAddress
}

func (p *sSettings) GetCapacity() uint64 {
	return p.FCapacity
}

func (p *sSettings) GetMaxConnects() uint64 {
	return p.FMaxConnects
}

func (p *sSettings) GetConnSettings() conn.ISettings {
	return p.FConnSettings
}
