package database

var (
	_ ISettings = &sSettings{}
)

const (
	cPath      = "database.db"
	cSaltKey   = "go-peer/salt"
	cCipherKey = "cipher-key"
)

type SSettings sSettings
type sSettings struct {
	FPath      string
	FHashing   bool
	FSaltKey   []byte
	FCipherKey []byte
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FPath:      pSett.FPath,
		FHashing:   pSett.FHashing,
		FSaltKey:   pSett.FSaltKey,
		FCipherKey: pSett.FCipherKey,
	}).useDefaultValues()
}

func (p *sSettings) useDefaultValues() ISettings {
	if p.FPath == "" {
		p.FPath = cPath
	}
	if p.FSaltKey == nil {
		p.FSaltKey = []byte(cSaltKey)
	}
	if p.FCipherKey == nil {
		p.FCipherKey = []byte(cCipherKey)
	}
	return p
}

func (p *sSettings) GetPath() string {
	return p.FPath
}

func (p *sSettings) GetSaltKey() []byte {
	return p.FSaltKey
}

func (p *sSettings) GetHashing() bool {
	return p.FHashing
}

func (p *sSettings) GetCipherKey() []byte {
	return p.FCipherKey
}
