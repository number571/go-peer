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
	FMessageSizeBytes uint64
	FLimitVoidSize    uint64
	FWaitReadDeadline time.Duration
	FReadDeadline     time.Duration
	FWriteDeadline    time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FNetworkKey:       pSett.FNetworkKey,
		FMessageSizeBytes: pSett.FMessageSizeBytes,
		FLimitVoidSize:    pSett.FLimitVoidSize,
		FWaitReadDeadline: pSett.FWaitReadDeadline,
		FReadDeadline:     pSett.FReadDeadline,
		FWriteDeadline:    pSett.FWriteDeadline,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FMessageSizeBytes == 0 {
		panic(`p.FMessageSizeBytes == 0`)
	}
	if p.FWaitReadDeadline == 0 {
		panic(`p.FWaitReadDeadline == 0`)
	}
	if p.FReadDeadline == 0 {
		panic(`p.FReadDeadline == 0`)
	}
	if p.FWriteDeadline == 0 {
		panic(`p.FWriteDeadline == 0`)
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

func (p *sSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}

func (p *sSettings) GetLimitVoidSize() uint64 {
	return p.FLimitVoidSize
}

func (p *sSettings) GetWaitReadDeadline() time.Duration {
	return p.FWaitReadDeadline
}

func (p *sSettings) GetReadDeadline() time.Duration {
	return p.FReadDeadline
}

func (p *sSettings) GetWriteDeadline() time.Duration {
	return p.FWriteDeadline
}
