package anonymity

import "time"

var (
	_ ISettings = &sSettings{}
)

const (
	cMaskNetwork = 0x676F2D7065657201
	cTimeWait    = time.Minute
)

type SSettings sSettings
type sSettings struct {
	FRetryEnqueue uint64
	FMaskNetwork  uint64
	FTimeWait     time.Duration
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FRetryEnqueue: sett.FRetryEnqueue,
		FMaskNetwork:  sett.FMaskNetwork,
		FTimeWait:     sett.FTimeWait,
	}).useDefaultValue()
}

func (s *sSettings) useDefaultValue() ISettings {
	if s.FMaskNetwork == 0 {
		s.FMaskNetwork = cMaskNetwork
	}
	if s.FTimeWait == 0 {
		s.FTimeWait = cTimeWait
	}
	return s
}

func (s *sSettings) GetTimeWait() time.Duration {
	return s.FTimeWait
}

func (s *sSettings) GetMaskNetwork() uint64 {
	return s.FMaskNetwork
}

func (s *sSettings) GetRetryEnqueue() uint64 {
	return s.FRetryEnqueue
}
