package client

const (
	cWorkSize    = 10
	cMessageSize = (1 << 20)
)

type SSettings struct {
	FWorkSize    uint64
	FMessageSize uint64
}

func NewSettings(sett *SSettings) ISettings {
	return (&SSettings{
		FWorkSize:    sett.FWorkSize,
		FMessageSize: sett.FMessageSize,
	}).useDefaultValues()
}

func (s *SSettings) useDefaultValues() ISettings {
	if s.FWorkSize == 0 {
		s.FWorkSize = cWorkSize
	}
	if s.FMessageSize == 0 {
		s.FMessageSize = cMessageSize
	}
	return s
}

func (s *SSettings) GetWorkSize() uint64 {
	return s.FWorkSize
}

func (s *SSettings) GetMessageSize() uint64 {
	return s.FMessageSize
}