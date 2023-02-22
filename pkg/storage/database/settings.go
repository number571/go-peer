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

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FPath:      sett.FPath,
		FHashing:   sett.FHashing,
		FSaltKey:   sett.FSaltKey,
		FCipherKey: sett.FCipherKey,
	}).useDefaultValues()
}

func (s *sSettings) useDefaultValues() ISettings {
	if s.FPath == "" {
		s.FPath = cPath
	}
	if s.FSaltKey == nil {
		s.FSaltKey = []byte(cSaltKey)
	}
	if s.FCipherKey == nil {
		s.FCipherKey = []byte(cCipherKey)
	}
	return s
}

func (s *sSettings) GetPath() string {
	return s.FPath
}

func (s *sSettings) GetSaltKey() []byte {
	return s.FSaltKey
}

func (s *sSettings) GetHashing() bool {
	return s.FHashing
}

func (s *sSettings) GetCipherKey() []byte {
	return s.FCipherKey
}
