package queue

import (
	"time"
)

const (
	cCapacity     = (1 << 5)
	cPullCapacity = (1 << 5)
	cDuration     = time.Second
)

type SSettings struct {
	FCapacity     uint64
	FPullCapacity uint64
	FDuration     time.Duration
}

func NewSettings(sett *SSettings) ISettings {
	return (&SSettings{
		FCapacity:     sett.FCapacity,
		FPullCapacity: sett.FPullCapacity,
		FDuration:     sett.FDuration,
	}).useDefaultValues()
}

func (s *SSettings) useDefaultValues() ISettings {
	if s.FCapacity == 0 {
		s.FCapacity = cCapacity
	}
	if s.FPullCapacity == 0 {
		s.FPullCapacity = cPullCapacity
	}
	if s.FDuration == 0 {
		s.FDuration = cDuration
	}
	return s
}

func (s *SSettings) GetCapacity() uint64 {
	return s.FCapacity
}

func (s *SSettings) GetPullCapacity() uint64 {
	return s.FPullCapacity
}

func (s *SSettings) GetDuration() time.Duration {
	return s.FDuration
}
