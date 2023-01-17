package database

var (
	_ ISettings = &sSettings{}
)

const (
	cWorkSize      = 10
	cMessageSize   = (1 << 20)
	cLimitMessages = (1 << 10)
)

type SSettings sSettings
type sSettings struct {
	FLimitMessages uint64
	FMessageSize   uint64
	FWorkSize      uint64
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FWorkSize:      sett.FWorkSize,
		FMessageSize:   sett.FMessageSize,
		FLimitMessages: sett.FLimitMessages,
	}).useDefaultValues()
}

func (s *sSettings) useDefaultValues() ISettings {
	if s.FWorkSize == 0 {
		s.FWorkSize = cWorkSize
	}
	if s.FMessageSize == 0 {
		s.FMessageSize = cMessageSize
	}
	if s.FLimitMessages == 0 {
		s.FLimitMessages = cLimitMessages
	}
	return s
}

func (s *sSettings) GetLimitMessages() uint64 {
	return s.FLimitMessages
}

func (s *sSettings) GetMessageSize() uint64 {
	return s.FMessageSize
}

func (s *sSettings) GetWorkSize() uint64 {
	return s.FWorkSize
}
