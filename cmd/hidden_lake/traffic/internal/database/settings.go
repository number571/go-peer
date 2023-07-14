package database

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FPath             string
	FMessageSizeBytes uint64
	FWorkSizeBits     uint64
	FCapacity         uint64
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FPath:             pSett.FPath,
		FWorkSizeBits:     pSett.FWorkSizeBits,
		FMessageSizeBytes: pSett.FMessageSizeBytes,
		FCapacity:         pSett.FCapacity,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FPath == "" {
		panic(`p.FPath == ""`)
	}
	if p.FWorkSizeBits == 0 {
		panic(`p.FWorkSizeBits == 0`)
	}
	if p.FMessageSizeBytes == 0 {
		panic(`p.FMessageSizeBytes == 0`)
	}
	// if capacity=0 -> then storage=false
	return p
}

func (p *sSettings) GetPath() string {
	return p.FPath
}

func (s *sSettings) GetCapacity() uint64 {
	return s.FCapacity
}

func (p *sSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}
