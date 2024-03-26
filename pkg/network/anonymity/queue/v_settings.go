package queue

var (
	_ ISettings = &sSettings{}
)

type SVSettings sVSettings
type sVSettings struct {
	FNetworkKey string
}

func NewVSettings(pSett *SVSettings) IVSettings {
	return (&sVSettings{
		FNetworkKey: pSett.FNetworkKey,
	}).mustNotNull()
}

func (p *sVSettings) mustNotNull() IVSettings {
	// p.FNetworkKey can be = 0
	return p
}

func (p *sVSettings) GetNetworkKey() string {
	return p.FNetworkKey
}
