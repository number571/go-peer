package anonymity

import "time"

const (
	cTimeWait = time.Minute
)

type SSettings struct {
	FRetryEnqueue uint64
	FTimeWait     time.Duration
}

func NewSettings(sett *SSettings) ISettings {
	return (&SSettings{
		FRetryEnqueue: sett.FRetryEnqueue,
		FTimeWait:     sett.FTimeWait,
	}).useDefaultValue()
}

func (s *SSettings) useDefaultValue() ISettings {
	if s.FTimeWait == 0 {
		s.FTimeWait = cTimeWait
	}
	return s
}

func (s *SSettings) GetTimeWait() time.Duration {
	return s.FTimeWait
}

func (s *SSettings) GetRetryEnqueue() uint64 {
	return s.FRetryEnqueue
}
