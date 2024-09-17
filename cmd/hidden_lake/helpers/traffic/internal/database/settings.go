package database

import "time"

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FPath             string
	FNetworkKey       string
	FWorkSizeBits     uint64
	FMessagesCapacity uint64
	FTimestampWindow  time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FPath:             pSett.FPath,
		FNetworkKey:       pSett.FNetworkKey,
		FWorkSizeBits:     pSett.FWorkSizeBits,
		FMessagesCapacity: pSett.FMessagesCapacity,
		FTimestampWindow:  pSett.FTimestampWindow,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FPath == "" {
		panic(`p.FPath == ""`)
	}
	if p.FMessagesCapacity == 0 {
		panic("p.FMessagesCapacity == 0")
	}
	return p
}

func (p *sSettings) GetPath() string {
	return p.FPath
}

func (p *sSettings) GetMessagesCapacity() uint64 {
	return p.FMessagesCapacity
}

func (p *sSettings) GetTimestampWindow() time.Duration {
	return p.FTimestampWindow
}

func (p *sSettings) GetNetworkKey() string {
	return p.FNetworkKey
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}
