package conn

var (
	_ IVSettings = &sVSettings{}
)

type SVSettings sVSettings
type sVSettings struct {
	FNetworkKey string
}

func NewVSettings(pVSett *SVSettings) IVSettings {
	return (&sVSettings{
		FNetworkKey: pVSett.FNetworkKey,
	}).mustNotNull()
}

func (p *sVSettings) mustNotNull() IVSettings {
	return p
}

func (p *sVSettings) GetNetworkKey() string {
	return p.FNetworkKey
}
