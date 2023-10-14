package message

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FWorkSizeBits     uint64
	FMessageSizeBytes uint64
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FWorkSizeBits:     pSett.FWorkSizeBits,
		FMessageSizeBytes: pSett.FMessageSizeBytes,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FMessageSizeBytes == 0 {
		panic(`p.FMessageSizeBytes == 0`)
	}
	return p
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *sSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}
