package conn_keeper

import (
	"errors"
	"sync"
	"time"

	"github.com/number571/go-peer/modules/network"
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

func (connKeeper *sConnKeeper) Settings() ISettings {
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
			case <-time.After(connKeeper.Settings().GetDuration()):
				connKeeper.tryConnectToAll()
			}
		}
	}()

	return nil
}

func (connKeeper *sConnKeeper) Close() error {
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
	for _, address := range connKeeper.Settings().GetConnections() {
		for addr := range connKeeper.fNode.Connections() {
			if addr == address {
				continue NEXT
			}
		}
		connKeeper.fNode.Connect(address)
	}
}