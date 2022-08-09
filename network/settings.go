package network

import "time"

type sSettings struct {
	fRetryNum    uint64
	fCapacity    uint64
	fSizePack    uint64
	fMaxConns    uint64
	fMaxMessages uint64
	fTimeWait    time.Duration
}

func NewSettings(sizePack, retryNum, capacity, maxConns, maxMessages uint64, timeWait time.Duration) ISettings {
	return &sSettings{
		fRetryNum:    retryNum,
		fCapacity:    capacity,
		fSizePack:    sizePack,
		fMaxConns:    maxConns,
		fMaxMessages: maxMessages,
		fTimeWait:    timeWait,
	}
}

func (s *sSettings) GetRetryNum() uint64 {
	return s.fCapacity
}

func (s *sSettings) GetCapacity() uint64 {
	return s.fCapacity
}

func (s *sSettings) GetPackageSize() uint64 {
	return s.fSizePack
}

func (s *sSettings) GetMaxConnects() uint64 {
	return s.fMaxConns
}

func (s *sSettings) GetMaxMessages() uint64 {
	return s.fMaxMessages
}

func (s *sSettings) GetTimeWait() time.Duration {
	return s.fTimeWait
}
