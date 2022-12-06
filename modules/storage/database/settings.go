package database

var (
	_ ISettings = &sSettings{}
)

const (
	cCipherKey = "cipher-key"
)

type SSettings sSettings
type sSettings struct {
	FHashing   bool
	FCipherKey []byte
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FHashing:   sett.FHashing,
		FCipherKey: sett.FCipherKey,
	}).useDefaultValues()
}

func (s *sSettings) useDefaultValues() ISettings {
	if s.FCipherKey == nil {
		s.FCipherKey = []byte(cCipherKey)
	}
	return s
}

func (s *sSettings) GetHashing() bool {
	return s.FHashing
}

func (s *sSettings) GetCipherKey() []byte {
	return s.FCipherKey
}
