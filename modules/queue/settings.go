package queue

import (
	"time"
)

const (
	cCapacity     = (1 << 5)
	cPullCapacity = (1 << 5)
	cDuration     = time.Second
)

type SSettings sSettings
type sSettings struct {
	FCapacity     uint64
	FPullCapacity uint64
	FDuration     time.Duration
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FCapacity:     sett.FCapacity,
		FPullCapacity: sett.FPullCapacity,
		FDuration:     sett.FDuration,
	}).useDefaultValues()
}

func (s *sSettings) useDefaultValues() ISettings {
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

func (s *sSettings) GetCapacity() uint64 {
	return s.FCapacity
}

func (s *sSettings) GetPullCapacity() uint64 {
	return s.FPullCapacity
}

func (s *sSettings) GetDuration() time.Duration {
	return s.FDuration
}
