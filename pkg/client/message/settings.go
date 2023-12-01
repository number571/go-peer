package message

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FKeySizeBits      uint64
	FMessageSizeBytes uint64
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FKeySizeBits:      pSett.FKeySizeBits,
		FMessageSizeBytes: pSett.FMessageSizeBytes,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FKeySizeBits == 0 {
		panic(`p.FKeySizeBits == 0`)
	}
	if p.FMessageSizeBytes == 0 {
		panic(`p.FMessageSizeBytes == 0`)
	}
	return p
}

func (p *sSettings) GetKeySizeBits() uint64 {
	return p.FKeySizeBits
}

func (p *sSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}
