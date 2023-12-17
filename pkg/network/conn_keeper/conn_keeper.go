package conn_keeper

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/state"
)

var (
	_ IConnKeeper = &sConnKeeper{}
)

type sConnKeeper struct {
	fState    state.IState
	fNode     network.INode
	fSettings ISettings
}

func NewConnKeeper(pSett ISettings, pNode network.INode) IConnKeeper {
	return &sConnKeeper{
		fState:    state.NewBoolState(),
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
	if err := p.fState.Enable(nil); err != nil {
		return fmt.Errorf("conn keeper running error: %w", err)
	}
	defer func() {
		if err := p.fState.Disable(nil); err != nil {
			panic(err)
		}
	}()

	for {
		p.tryConnectToAll()
		select {
		case <-pCtx.Done():
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
