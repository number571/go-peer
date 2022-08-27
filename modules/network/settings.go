package network

import "time"

const (
	cRetryNum    = 2
	cCapacity    = (1 << 10)
	cMessageSize = (1 << 20)
	cMaxConns    = (1 << 6)
	cMaxMessages = (1 << 4)
	cTimeWait    = time.Minute
)

type SSettings struct {
	FRetryNum    uint64
	FCapacity    uint64
	FMessageSize uint64
	FMaxConns    uint64
	FMaxMessages uint64
	FTimeWait    time.Duration
}

func NewSettings(sett *SSettings) ISettings {
	return (&SSettings{
		FRetryNum:    sett.FRetryNum,
		FCapacity:    sett.FCapacity,
		FMessageSize: sett.FMessageSize,
		FMaxConns:    sett.FMaxConns,
		FMaxMessages: sett.FMaxMessages,
		FTimeWait:    sett.FTimeWait,
	}).useDefaultValues()
}

func (s *SSettings) useDefaultValues() ISettings {
	if s.FRetryNum == 0 {
		s.FRetryNum = cRetryNum
	}
	if s.FCapacity == 0 {
		s.FCapacity = cCapacity
	}
	if s.FMessageSize == 0 {
		s.FMessageSize = cMessageSize
	}
	if s.FMaxConns == 0 {
		s.FMaxConns = cMaxConns
	}
	if s.FMaxMessages == 0 {
		s.FMaxMessages = cMaxMessages
	}
	if s.FTimeWait == 0 {
		s.FTimeWait = cTimeWait
	}
	return s
}

func (s *SSettings) GetRetryNum() uint64 {
	return s.FCapacity
}

func (s *SSettings) GetCapacity() uint64 {
	return s.FCapacity
}

func (s *SSettings) GetMessageSize() uint64 {
	return s.FMessageSize
}

func (s *SSettings) GetMaxConnects() uint64 {
	return s.FMaxConns
}

func (s *SSettings) GetMaxMessages() uint64 {
	return s.FMaxMessages
}

func (s *SSettings) GetTimeWait() time.Duration {
	return s.FTimeWait
}
