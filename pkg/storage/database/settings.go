package database

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FPath      string
	FHashing   bool
	FCipherKey []byte
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FPath:      pSett.FPath,
		FHashing:   pSett.FHashing,
		FCipherKey: pSett.FCipherKey,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FPath == "" {
		panic(`p.FPath == ""`)
	}
	if p.FCipherKey == nil {
		panic(`p.FCipherKey == nil`)
	}
	return p
}

func (p *sSettings) GetPath() string {
	return p.FPath
}

func (p *sSettings) GetHashing() bool {
	return p.FHashing
}

func (p *sSettings) GetCipherKey() []byte {
	return p.FCipherKey
}
