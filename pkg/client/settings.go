package client

var (
	_ ISettings = &sSettings{}
)

const (
	cWorkSize    = 10
	cMessageSize = (1 << 20)
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
	}).useDefaultValues()
}

func (p *sSettings) useDefaultValues() ISettings {
	if p.FWorkSize == 0 {
		p.FWorkSize = cWorkSize
	}
	if p.FMessageSize == 0 {
		p.FMessageSize = cMessageSize
	}
	return p
}

func (p *sSettings) GetWorkSize() uint64 {
	return p.FWorkSize
}

func (p *sSettings) GetMessageSize() uint64 {
	return p.FMessageSize
}
