package database

const (
	cPath      = "database.db"
	cCipherKey = "cipher-key"
)

type SSettings struct {
	FPath      string
	FHashing   bool
	FCipherKey []byte
}

func NewSettings(sett *SSettings) ISettings {
	return (&SSettings{
		FPath:      sett.FPath,
		FHashing:   sett.FHashing,
		FCipherKey: sett.FCipherKey,
	}).useDefaultValues()
}

func (s *SSettings) useDefaultValues() ISettings {
	if s.FPath == "" {
		s.FPath = cPath
	}
	if s.FCipherKey == nil {
		s.FCipherKey = []byte(cCipherKey)
	}
	return s
}

func (s *SSettings) GetPath() string {
	return s.FPath
}

func (s *SSettings) GetHashing() bool {
	return s.FHashing
}

func (s *SSettings) GetCipherKey() []byte {
	return s.FCipherKey
}
