package conn

import (
	"time"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FNetworkKey            string
	FWorkSizeBits          uint64
	FLimitMessageSizeBytes uint64
	FTimestampWindow       time.Duration
	FWaitReadTimeout       time.Duration
	FDialTimeout           time.Duration
	FReadTimeout           time.Duration
	FWriteTimeout          time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FNetworkKey:            pSett.FNetworkKey,
		FWorkSizeBits:          pSett.FWorkSizeBits,
		FLimitMessageSizeBytes: pSett.FLimitMessageSizeBytes,
		FTimestampWindow:       pSett.FTimestampWindow,
		FWaitReadTimeout:       pSett.FWaitReadTimeout,
		FDialTimeout:           pSett.FDialTimeout,
		FReadTimeout:           pSett.FReadTimeout,
		FWriteTimeout:          pSett.FWriteTimeout,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FLimitMessageSizeBytes == 0 {
		panic(`p.FLimitMessageSizeBytes == 0`)
	}
	if p.FWaitReadTimeout == 0 {
		panic(`p.FWaitReadTimeout == 0`)
	}
	if p.FDialTimeout == 0 {
		panic(`p.FDialTimeout == 0`)
	}
	if p.FReadTimeout == 0 {
		panic(`p.FReadTimeout == 0`)
	}
	if p.FWriteTimeout == 0 {
		panic(`p.FWriteTimeout == 0`)
	}
	return p
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *sSettings) GetLimitMessageSizeBytes() uint64 {
	return p.FLimitMessageSizeBytes
}

func (p *sSettings) GetTimestampWindow() time.Duration {
	return p.FTimestampWindow
}

func (p *sSettings) GetWaitReadTimeout() time.Duration {
	return p.FWaitReadTimeout
}

func (p *sSettings) GetDialTimeout() time.Duration {
	return p.FDialTimeout
}

func (p *sSettings) GetReadTimeout() time.Duration {
	return p.FReadTimeout
}

func (p *sSettings) GetWriteTimeout() time.Duration {
	return p.FWriteTimeout
}
