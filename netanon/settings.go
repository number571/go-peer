package netanon

import "time"

type sSettings struct {
	fResponsePeriod uint64
	fRetryEnqueue   uint64
	fTimeWait       time.Duration
}

func NewSettings(respPeriod uint64, retryNum uint64, timeWait time.Duration) ISettings {
	return &sSettings{
		fResponsePeriod: respPeriod,
		fRetryEnqueue:   retryNum,
		fTimeWait:       timeWait,
	}
}

func (s *sSettings) GetResponsePeriod() uint64 {
	return s.fResponsePeriod
}

func (s *sSettings) GetRetryEnqueue() uint64 {
	return s.fRetryEnqueue
}

func (s *sSettings) GetTimeWait() time.Duration {
	return s.fTimeWait
}
