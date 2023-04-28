package transaction

var (
	_ ISettings = &sSettings{}
)

const (
	cMaxSize = (1 << 10)
)

type SSettings sSettings
type sSettings struct {
	FMaxSize uint64
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FMaxSize: sett.FMaxSize,
	}).useDefaultValues()
}

func (s *sSettings) useDefaultValues() ISettings {
	if s.FMaxSize == 0 {
		s.FMaxSize = cMaxSize
	}
	return s
}

func (sett *sSettings) GetMaxSize() uint64 {
	return sett.FMaxSize
}
