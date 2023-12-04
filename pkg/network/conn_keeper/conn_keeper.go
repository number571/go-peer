package conn_keeper

import (
	"errors"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/network"
)

var (
	_ IConnKeeper = &sConnKeeper{}
)

type sConnKeeper struct {
	fIsRun    bool
	fMutex    sync.RWMutex
	fSignal   chan struct{}
	fNode     network.INode
	fSettings ISettings
}

func NewConnKeeper(pSett ISettings, pNode network.INode) IConnKeeper {
	return &sConnKeeper{
		fNode:     pNode,
		fSettings: pSett,
	}
}

func (p *sConnKeeper) GetNetworkNode() network.INode {
	return p.fNode
}

func (p *sConnKeeper) GetSettings() ISettings {
	return p.fSettings
}

func (p *sConnKeeper) Run() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fIsRun {
		return errors.New("conn keeper already started")
	}
	p.fIsRun = true

	p.fSignal = make(chan struct{})
	p.tryConnectToAll()

	go func() {
		for {
			select {
			case <-p.readSignal():
				return
			case <-time.After(p.fSettings.GetDuration()):
				p.tryConnectToAll()
			}
		}
	}()

	return nil
}

func (p *sConnKeeper) Stop() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if !p.fIsRun {
		return errors.New("conn keeper already closed or not started")
	}
	p.fIsRun = false

	close(p.fSignal)
	return nil
}

func (p *sConnKeeper) tryConnectToAll() {
NEXT:
	for _, address := range p.fSettings.GetConnections() {
		mapConns := p.fNode.GetConnections()
		if _, ok := mapConns[address]; ok {
			continue NEXT
		}
		p.fNode.AddConnection(address)
	}
}

func (p *sConnKeeper) readSignal() <-chan struct{} {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return p.fSignal
}
