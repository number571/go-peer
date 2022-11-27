package client

const (
	cWorkSize    = 10
	cMessageSize = (1 << 20)
)

type SSettings sSettings
type sSettings struct {
	FWorkSize    uint64
	FMessageSize uint64
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FWorkSize:    sett.FWorkSize,
		FMessageSize: sett.FMessageSize,
	}).useDefaultValues()
}

func (s *sSettings) useDefaultValues() ISettings {
	if s.FWorkSize == 0 {
		s.FWorkSize = cWorkSize
	}
	if s.FMessageSize == 0 {
		s.FMessageSize = cMessageSize
	}
	return s
}

func (s *sSettings) GetWorkSize() uint64 {
	return s.FWorkSize
}

func (s *sSettings) GetMessageSize() uint64 {
	return s.FMessageSize
}
