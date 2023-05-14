package message

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FWorkSize    uint64
	FMessageSize uint64
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FWorkSize:    pSett.FWorkSize,
		FMessageSize: pSett.FMessageSize,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FWorkSize == 0 {
		panic(`p.FWorkSize == 0`)
	}
	if p.FMessageSize == 0 {
		panic(`p.FMessageSize == 0`)
	}
	return p
}

func (p *sSettings) GetWorkSize() uint64 {
	return p.FWorkSize
}

func (p *sSettings) GetMessageSize() uint64 {
	return p.FMessageSize
}
