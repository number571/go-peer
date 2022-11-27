package network

import "github.com/number571/go-peer/modules/network/conn"

const (
	cCapacity    = (1 << 10)
	cMaxConnects = (1 << 6)
)

type SSettings sSettings
type sSettings struct {
	FCapacity     uint64
	FMaxConnects  uint64
	FConnSettings conn.ISettings
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FCapacity:     sett.FCapacity,
		FMaxConnects:  sett.FMaxConnects,
		FConnSettings: sett.FConnSettings,
	}).useDefaultValues()
}

func (s *sSettings) useDefaultValues() ISettings {
	if s.FCapacity == 0 {
		s.FCapacity = cCapacity
	}
	if s.FMaxConnects == 0 {
		s.FMaxConnects = cMaxConnects
	}
	if s.FConnSettings == nil {
		s.FConnSettings = conn.NewSettings(&conn.SSettings{})
	}
	return s
}

func (s *sSettings) GetCapacity() uint64 {
	return s.FCapacity
}

func (s *sSettings) GetMaxConnects() uint64 {
	return s.FMaxConnects
}

func (s *sSettings) GetConnSettings() conn.ISettings {
	return s.FConnSettings
}
