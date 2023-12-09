package database

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FPath             string
	FNetworkKey       string
	FWorkSizeBits     uint64
	FHashesWindow     uint64
	FMessagesCapacity uint64
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FPath:             pSett.FPath,
		FNetworkKey:       pSett.FNetworkKey,
		FWorkSizeBits:     pSett.FWorkSizeBits,
		FHashesWindow:     pSett.FHashesWindow,
		FMessagesCapacity: pSett.FMessagesCapacity,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FPath == "" {
		panic(`p.FPath == ""`)
	}
	if p.FHashesWindow == 0 {
		panic("p.FHashesWindow == 0")
	}
	if p.FMessagesCapacity == 0 {
		panic("p.FMessagesCapacity == 0")
	}
	return p
}

func (p *sSettings) GetPath() string {
	return p.FPath
}

func (p *sSettings) GetHashesWindow() uint64 {
	return p.FHashesWindow
}

func (s *sSettings) GetMessagesCapacity() uint64 {
	return s.FMessagesCapacity
}

func (s *sSettings) GetNetworkKey() string {
	return s.FNetworkKey
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}
