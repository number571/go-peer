package message

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FEncKeySizeBytes  uint64
	FMessageSizeBytes uint64
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FEncKeySizeBytes:  pSett.FEncKeySizeBytes,
		FMessageSizeBytes: pSett.FMessageSizeBytes,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FEncKeySizeBytes == 0 {
		panic(`p.FEncKeySizeBytes == 0`)
	}
	if p.FMessageSizeBytes == 0 {
		panic(`p.FMessageSizeBytes == 0`)
	}
	return p
}

func (p *sSettings) GetEncKeySizeBytes() uint64 {
	return p.FEncKeySizeBytes
}

func (p *sSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}
