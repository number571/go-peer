package conn

import (
	"sync"
	"time"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	fMutex sync.Mutex

	FNetworkKey       string
	FWorkSizeBits     uint64
	FMessageSizeBytes uint64
	FLimitVoidSize    uint64
	FWaitReadTimeout  time.Duration
	FDialTimeout      time.Duration
	FReadTimeout      time.Duration
	FWriteTimeout     time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FNetworkKey:       pSett.FNetworkKey,
		FWorkSizeBits:     pSett.FWorkSizeBits,
		FMessageSizeBytes: pSett.FMessageSizeBytes,
		FLimitVoidSize:    pSett.FLimitVoidSize,
		FWaitReadTimeout:  pSett.FWaitReadTimeout,
		FDialTimeout:      pSett.FDialTimeout,
		FReadTimeout:      pSett.FReadTimeout,
		FWriteTimeout:     pSett.FWriteTimeout,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FMessageSizeBytes == 0 {
		panic(`p.FMessageSizeBytes == 0`)
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
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.FNetworkKey
}

func (p *sSettings) SetNetworkKey(pNetworkKey string) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.FNetworkKey = pNetworkKey
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *sSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}

func (p *sSettings) GetLimitVoidSize() uint64 {
	return p.FLimitVoidSize
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
