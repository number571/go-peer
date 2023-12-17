package conn_keeper

import (
	"context"
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
	fMutex    sync.Mutex
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

func (p *sConnKeeper) Run(pCtx context.Context) error {
	err := func() error {
		p.fMutex.Lock()
		defer p.fMutex.Unlock()

		if p.fIsRun {
			return errors.New("conn keeper already running")
		}

		p.fIsRun = true
		return nil
	}()
	if err != nil {
		return err
	}

	for {
		p.tryConnectToAll()
		select {
		case <-pCtx.Done():
			p.fMutex.Lock()
			p.fIsRun = false
			p.fMutex.Unlock()
			return nil
		case <-time.After(p.fSettings.GetDuration()):
			// next iter
		}
	}
}

func (p *sConnKeeper) tryConnectToAll() {
	connList := p.fSettings.GetConnections()

	wg := sync.WaitGroup{}
	wg.Add(len(connList))

	for _, addr := range connList {
		go func(addr string) {
			defer wg.Done()

			mapConns := p.fNode.GetConnections()
			if _, ok := mapConns[addr]; ok {
				return
			}

			p.fNode.AddConnection(addr)
		}(addr)
	}

	wg.Wait()
}
