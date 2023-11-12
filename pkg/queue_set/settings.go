package queue_set

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FCapacity uint64
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FCapacity: pSett.FCapacity,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FCapacity == 0 {
		panic(`p.FCapacity == 0`)
	}
	return p
}

func (p *sSettings) GetCapacity() uint64 {
	return p.FCapacity
}
