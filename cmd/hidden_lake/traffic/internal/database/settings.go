package database

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FPath        string
	FMessageSize uint64
	FWorkSize    uint64
	FCapacity    uint64
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FPath:        pSett.FPath,
		FWorkSize:    pSett.FWorkSize,
		FMessageSize: pSett.FMessageSize,
		FCapacity:    pSett.FCapacity,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FPath == "" {
		panic(`p.FPath == ""`)
	}
	if p.FWorkSize == 0 {
		panic(`p.FWorkSize == 0`)
	}
	if p.FMessageSize == 0 {
		panic(`p.FMessageSize == 0`)
	}
	if p.FCapacity == 0 {
		panic(`p.FLimitMessages == 0`)
	}
	return p
}

func (p *sSettings) GetPath() string {
	return p.FPath
}

func (s *sSettings) GetCapacity() uint64 {
	return s.FCapacity
}

func (p *sSettings) GetMessageSize() uint64 {
	return p.FMessageSize
}

func (p *sSettings) GetWorkSize() uint64 {
	return p.FWorkSize
}
