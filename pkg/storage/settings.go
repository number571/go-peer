package storage

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FPath     string
	FWorkSize uint64
	FPassword string
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FPath:     pSett.FPath,
		FWorkSize: pSett.FWorkSize,
		FPassword: pSett.FPassword,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FPath == "" {
		panic(`p.FPath == ""`)
	}
	if p.FWorkSize == 0 {
		panic(`p.FWorkSize == 0`)
	}
	if p.FPassword == "" {
		panic(`p.FPassword == ""`)
	}
	return p
}

func (p *sSettings) GetPath() string {
	return p.FPath
}

func (p *sSettings) GetWorkSize() uint64 {
	return p.FWorkSize
}

func (p *sSettings) GetPassword() string {
	return p.FPassword
}
