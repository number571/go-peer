package database

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FPath string
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FPath: pSett.FPath,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FPath == "" {
		panic(`p.FPath == ""`)
	}
	return p
}

func (p *sSettings) GetPath() string {
	return p.FPath
}
