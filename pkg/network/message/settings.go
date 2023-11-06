package message

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FWorkSizeBits uint64
	FNetworkKey   string
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FWorkSizeBits: pSett.FWorkSizeBits,
		FNetworkKey:   pSett.FNetworkKey,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	return p
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *sSettings) GetNetworkKey() string {
	return p.FNetworkKey
}
