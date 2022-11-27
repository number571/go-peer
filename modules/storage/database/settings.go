package database

const (
	cPath      = "database.db"
	cCipherKey = "cipher-key"
)

type SSettings sSettings
type sSettings struct {
	FPath      string
	FHashing   bool
	FCipherKey []byte
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FPath:      sett.FPath,
		FHashing:   sett.FHashing,
		FCipherKey: sett.FCipherKey,
	}).useDefaultValues()
}

func (s *sSettings) useDefaultValues() ISettings {
	if s.FPath == "" {
		s.FPath = cPath
	}
	if s.FCipherKey == nil {
		s.FCipherKey = []byte(cCipherKey)
	}
	return s
}

func (s *sSettings) GetPath() string {
	return s.FPath
}

func (s *sSettings) GetHashing() bool {
	return s.FHashing
}

func (s *sSettings) GetCipherKey() []byte {
	return s.FCipherKey
}
