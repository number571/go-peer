package storage

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FPath      string
	FWorkSize  uint64
	FCipherKey []byte
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FPath:      pSett.FPath,
		FWorkSize:  pSett.FWorkSize,
		FCipherKey: pSett.FCipherKey,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FPath == "" {
		panic(`p.FPath == ""`)
	}
	if p.FWorkSize == 0 {
		panic(`p.FWorkSize == 0`)
	}
	if p.FCipherKey == nil {
		panic(`p.FCipherKey == nil`)
	}
	return p
}

func (p *sSettings) GetPath() string {
	return p.FPath
}

func (p *sSettings) GetWorkSize() uint64 {
	return p.FWorkSize
}

func (p *sSettings) GetCipherKey() []byte {
	return p.FCipherKey
}
