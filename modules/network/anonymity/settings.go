package anonymity

import "time"

const (
	cTimeWait = time.Minute
)

type SSettings sSettings
type sSettings struct {
	FRetryEnqueue uint64
	FTimeWait     time.Duration
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FRetryEnqueue: sett.FRetryEnqueue,
		FTimeWait:     sett.FTimeWait,
	}).useDefaultValue()
}

func (s *sSettings) useDefaultValue() ISettings {
	if s.FTimeWait == 0 {
		s.FTimeWait = cTimeWait
	}
	return s
}

func (s *sSettings) GetTimeWait() time.Duration {
	return s.FTimeWait
}

func (s *sSettings) GetRetryEnqueue() uint64 {
	return s.FRetryEnqueue
}
