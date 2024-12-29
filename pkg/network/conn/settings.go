package conn

import (
	"time"

	"github.com/number571/go-peer/pkg/message/layer1"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FLimitMessageSizeBytes uint64
	FWaitReadTimeout       time.Duration
	FDialTimeout           time.Duration
	FReadTimeout           time.Duration
	FWriteTimeout          time.Duration
	FMessageSettings       layer1.ISettings
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FMessageSettings:       pSett.FMessageSettings,
		FLimitMessageSizeBytes: pSett.FLimitMessageSizeBytes,
		FWaitReadTimeout:       pSett.FWaitReadTimeout,
		FDialTimeout:           pSett.FDialTimeout,
		FReadTimeout:           pSett.FReadTimeout,
		FWriteTimeout:          pSett.FWriteTimeout,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FMessageSettings == nil {
		panic(`p.FMessageSettings == nil`)
	}
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

func (p *sSettings) GetMessageSettings() layer1.ISettings {
	return p.FMessageSettings
}

func (p *sSettings) GetLimitMessageSizeBytes() uint64 {
	return p.FLimitMessageSizeBytes
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
