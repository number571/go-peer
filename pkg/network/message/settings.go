package message

var (
	_ IConstructSettings = &sConstructSettings{}
	_ ISettings          = &sSettings{}
)

type SConstructSettings sConstructSettings
type sConstructSettings struct {
	FSettings             ISettings
	FParallel             uint64
	FRandMessageSizeBytes uint64
}

type SSettings sSettings
type sSettings struct {
	FWorkSizeBits uint64
	FNetworkKey   string
}

func NewConstructSettings(pSett *SConstructSettings) IConstructSettings {
	return (&sConstructSettings{
		FSettings: func() ISettings {
			if pSett.FSettings == nil {
				return NewSettings(&SSettings{})
			}
			return pSett.FSettings
		}(),
		FParallel:             pSett.FParallel,
		FRandMessageSizeBytes: pSett.FRandMessageSizeBytes,
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

func (p *sConstructSettings) GetNetworkKey() string {
	return p.FSettings.GetNetworkKey()
}

func (p *sConstructSettings) GetWorkSizeBits() uint64 {
	return p.FSettings.GetWorkSizeBits()
}

func (p *sConstructSettings) GetParallel() uint64 {
	return p.FParallel
}

func (p *sConstructSettings) GetRandMessageSizeBytes() uint64 {
	return p.FRandMessageSizeBytes
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
