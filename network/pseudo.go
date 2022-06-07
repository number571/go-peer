package network

import (
	"sync"
	"time"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/settings"
)

var (
	_ iPseudo = &sPseudo{}
)

type sPseudo struct {
	fMutex   sync.Mutex
	fNode    INode
	fEnabled bool
	fChannel chan struct{}
	fPrivKey crypto.IPrivKey
}

// Set state = bool.
func (psd *sPseudo) Switch(state bool) {
	psd.fMutex.Lock()
	defer psd.fMutex.Unlock()

	if psd.fEnabled == state {
		return
	}
	psd.fEnabled = state

	switch state {
	case true:
		psd.start()
	case false:
		psd.stop()
	}
}

// Get current state of online mode.
func (psd *sPseudo) Status() bool {
	psd.fMutex.Lock()
	defer psd.fMutex.Unlock()

	return psd.fEnabled
}

func (psd *sPseudo) Request(size int) iPseudo {
	psd.fMutex.Lock()
	defer psd.fMutex.Unlock()

	if !psd.fEnabled {
		return psd
	}

	psd.doRequest(size)
	return psd
}

func (psd *sPseudo) Sleep() iPseudo {
	psd.fMutex.Lock()
	defer psd.fMutex.Unlock()

	if !psd.fEnabled {
		return psd
	}

	node := psd.fNode.(*sNode)
	wtime := node.fClient.Settings().Get(settings.CTimePrsp)
	time.Sleep(time.Millisecond * calcRandTime(wtime))
	return psd
}

// Get pseudo public key.
func (psd *sPseudo) PubKey() crypto.IPubKey {
	psd.fMutex.Lock()
	defer psd.fMutex.Unlock()

	return psd.fPrivKey.PubKey()
}

// Get pseudo private key.
func (psd *sPseudo) PrivKey() crypto.IPrivKey {
	psd.fMutex.Lock()
	defer psd.fMutex.Unlock()

	return psd.fPrivKey
}

func (psd *sPseudo) start() {
	go func(psd *sPseudo) {
		sett := psd.fNode.(*sNode).fClient.Settings()
		for {
			psd.doRequest(16)
			select {
			case <-psd.fChannel:
				return
			case <-time.After(time.Second * time.Duration(
				sett.Get(settings.CTimePreq),
			)):
				continue
			}
		}
	}(psd)
}

func (psd *sPseudo) stop() {
	psd.fMutex.Unlock()
	psd.fChannel <- struct{}{}
	psd.fMutex.Lock()
}

func (psd *sPseudo) doRequest(size int) {
	node := psd.fNode.(*sNode)
	rand := crypto.NewPRNG()
	pMsg, _ := node.fClient.Encrypt(
		local.NewRoute(psd.fPrivKey.PubKey()),
		local.NewMessage(
			rand.Bytes(16),
			rand.Bytes(calcRandSize(size)),
		),
	)
	ch := make(chan struct{})
	go func(node *sNode, ch chan struct{}) {
		node.send(pMsg)
		ch <- struct{}{}
	}(node, ch)
	<-ch
}

func calcRandSize(len int) uint64 {
	ulen := uint64(len)
	rand := crypto.NewPRNG()
	return ulen + rand.Uint64()%(10<<10) // +[0;10]KiB
}

func calcRandTime(seconds uint64) time.Duration {
	rand := crypto.NewPRNG()
	return time.Duration(rand.Uint64() % (seconds * 1000)) // random[0;S*1000]MS
}
