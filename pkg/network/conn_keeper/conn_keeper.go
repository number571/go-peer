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
	fMutex    sync.Mutex
	fSignal   chan struct{}
	fNode     network.INode
	fSettings ISettings
}

func NewConnKeeper(sett ISettings, node network.INode) IConnKeeper {
	return &sConnKeeper{
		fNode:     node,
		fSettings: sett,
	}
}

func (connKeeper *sConnKeeper) GetNetworkNode() network.INode {
	return connKeeper.fNode
}

func (connKeeper *sConnKeeper) GetSettings() ISettings {
	return connKeeper.fSettings
}

func (connKeeper *sConnKeeper) Run() error {
	connKeeper.fMutex.Lock()
	defer connKeeper.fMutex.Unlock()

	if connKeeper.fIsRun {
		return errors.New("conn keeper already started")
	}
	connKeeper.fIsRun = true

	connKeeper.fSignal = make(chan struct{})
	connKeeper.tryConnectToAll()

	go func() {
		for {
			select {
			case <-connKeeper.fSignal:
				return
			case <-time.After(connKeeper.GetSettings().GetDuration()):
				connKeeper.tryConnectToAll()
			}
		}
	}()

	return nil
}

func (connKeeper *sConnKeeper) Stop() error {
	connKeeper.fMutex.Lock()
	defer connKeeper.fMutex.Unlock()

	if !connKeeper.fIsRun {
		return errors.New("conn keeper already closed or not started")
	}
	connKeeper.fIsRun = false

	close(connKeeper.fSignal)
	return nil
}

func (connKeeper *sConnKeeper) tryConnectToAll() {
NEXT:
	for _, address := range connKeeper.GetSettings().GetConnections() {
		for addr := range connKeeper.fNode.GetConnections() {
			if addr == address {
				continue NEXT
			}
		}
		connKeeper.fNode.AddConnect(address)
	}
}
