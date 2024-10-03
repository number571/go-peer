package conn

import (
	"time"

	net_message "github.com/number571/go-peer/pkg/network/message"
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
	FMessageSettings       net_message.ISettings
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FMessageSettings: func() net_message.ISettings {
			if pSett.FMessageSettings == nil {
				return net_message.NewSettings(&net_message.SSettings{})
			}
			return pSett.FMessageSettings
		}(),
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

func (p *sSettings) GetNetworkKey() string {
	return p.FMessageSettings.GetNetworkKey()
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FMessageSettings.GetWorkSizeBits()
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
