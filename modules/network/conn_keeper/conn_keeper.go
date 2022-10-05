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

type iState int

const (
	cIsInit iState = iota
	cIsRun
	cIsClose
)

type sConnKeeper struct {
	fMutex    sync.Mutex
	fState    iState
	fSignal   chan struct{}
	fNode     network.INode
	fSettings ISettings
}

func NewConnKeeper(sett ISettings, node network.INode) IConnKeeper {
	return &sConnKeeper{
		fState:    cIsInit,
		fSignal:   make(chan struct{}),
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

	if connKeeper.fState != cIsInit {
		return errors.New("conn keeper already started or closed")
	}
	connKeeper.fState = cIsRun

	go func() {
		for {
			select {
			case <-connKeeper.fSignal:
				connKeeper.fState = cIsClose
				return
			default:
				connKeeper.tryConnectToAll()
				time.Sleep(connKeeper.Settings().GetDuration())
			}
		}
	}()

	return nil
}

func (connKeeper *sConnKeeper) Close() error {
	connKeeper.fMutex.Lock()
	defer connKeeper.fMutex.Unlock()

	if connKeeper.fState != cIsRun {
		return errors.New("conn keeper already closed or not started")
	}

	connKeeper.fSignal <- struct{}{}
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
