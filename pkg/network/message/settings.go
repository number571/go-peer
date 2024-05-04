package message

var (
	_ IConstructSettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FWorkSizeBits       uint64
	FNetworkKey         string
	FParallel           uint64
	FLimitVoidSizeBytes uint64
}

func NewSettings(pSett *SSettings) IConstructSettings {
	return (&sSettings{
		FWorkSizeBits:       pSett.FWorkSizeBits,
		FNetworkKey:         pSett.FNetworkKey,
		FParallel:           pSett.FParallel,
		FLimitVoidSizeBytes: pSett.FLimitVoidSizeBytes,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() IConstructSettings {
	return p
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *sSettings) GetNetworkKey() string {
	return p.FNetworkKey
}

func (p *sSettings) GetParallel() uint64 {
	return p.FParallel
}

func (p *sSettings) GetLimitVoidSizeBytes() uint64 {
	return p.FLimitVoidSizeBytes
}
