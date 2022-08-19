package netanon

import "time"

type sSettings struct {
	fRetryEnqueue uint64
	fTimeWait     time.Duration
}

func NewSettings(retryNum uint64, timeWait time.Duration) ISettings {
	return &sSettings{
		fRetryEnqueue: retryNum,
		fTimeWait:     timeWait,
	}
}

func (s *sSettings) GetTimeWait() time.Duration {
	return s.fTimeWait
}

func (s *sSettings) GetRetryEnqueue() uint64 {
	return s.fRetryEnqueue
}
