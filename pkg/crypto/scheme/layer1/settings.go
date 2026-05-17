package layer1

var (
	_ IConstructSettings = &sConstructSettings{}
	_ ISettings          = &sSettings{}
)

type SConstructSettings sConstructSettings
type sConstructSettings struct {
	FSettings ISettings
	FParallel uint64
}

type SSettings sSettings
type sSettings struct {
	FWorkSizeBits uint64
	FNetworkKey   string
}

func NewConstructSettings(pSett *SConstructSettings) IConstructSettings {
	return (&sConstructSettings{
		FSettings: pSett.FSettings,
		FParallel: pSett.FParallel,
	}).mustNotNull()
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FWorkSizeBits: pSett.FWorkSizeBits,
		FNetworkKey:   pSett.FNetworkKey,
	}).mustNotNull()
}

func (p *sConstructSettings) mustNotNull() IConstructSettings {
	if p.FSettings == nil {
		panic(`p.FSettings == nil`)
	}
	return p
}

func (p *sConstructSettings) GetSettings() ISettings {
	return p.FSettings
}

func (p *sConstructSettings) GetParallel() uint64 {
	return p.FParallel
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
