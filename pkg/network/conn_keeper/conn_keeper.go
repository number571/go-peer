package conn_keeper

import (
	"context"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/utils"
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
		return utils.MergeErrors(ErrRunning, err)
	}
	defer func() {
		if err := p.fState.Disable(nil); err != nil {
			panic(err)
		}
	}()

	for {
		p.tryConnectToAll(pCtx)
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		case <-time.After(p.fSettings.GetDuration()):
			// next iter
		}
	}
}

func (p *sConnKeeper) tryConnectToAll(pCtx context.Context) {
	connList := p.fSettings.GetConnections()
	mapConns := p.fNode.GetConnections()

	wg := sync.WaitGroup{}
	wg.Add(len(connList))

	for _, addr := range connList {
		go func(addr string) {
			defer wg.Done()
			if _, ok := mapConns[addr]; ok {
				return
			}
			p.fNode.AddConnection(pCtx, addr)
		}(addr)
	}

	wg.Wait()
}
