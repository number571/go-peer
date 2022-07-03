package netanon

import (
	"sync"
	"time"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/crypto/random"
	"github.com/number571/go-peer/local/payload"
	"github.com/number571/go-peer/local/routing"
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
	fPrivKey asymmetric.IPrivKey
}

func newPseudo(node INode) iPseudo {
	return &sPseudo{
		fNode:    node,
		fChannel: make(chan struct{}),
		fPrivKey: asymmetric.NewRSAPrivKey(node.Client().PubKey().Size()),
	}
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

func (psd *sPseudo) request(size int) iPseudo {
	psd.fMutex.Lock()
	defer psd.fMutex.Unlock()

	if !psd.fEnabled {
		return psd
	}

	psd.doRequest(size)
	return psd
}

func (psd *sPseudo) sleep() iPseudo {
	psd.fMutex.Lock()
	defer psd.fMutex.Unlock()

	if !psd.fEnabled {
		return psd
	}

	wtime := psd.fNode.Client().Settings().Get(settings.CTimePrsp)
	time.Sleep(time.Millisecond * calcRandTime(wtime))
	return psd
}

// Get pseudo private key.
func (psd *sPseudo) privKey() asymmetric.IPrivKey {
	psd.fMutex.Lock()
	defer psd.fMutex.Unlock()

	return psd.fPrivKey
}

func (psd *sPseudo) start() {
	sett := psd.fNode.Client().Settings()
	go func(psd *sPseudo, sett settings.ISettings) {
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
	}(psd, sett)
}

func (psd *sPseudo) stop() {
	psd.fChannel <- struct{}{}
}

func (psd *sPseudo) doRequest(size int) {
	rand := random.NewStdPRNG()
	psd.fNode.Broadcast(psd.fNode.Client().Encrypt(
		routing.NewRoute(psd.fPrivKey.PubKey()),
		payload.NewPayload(
			psd.fNode.Client().Settings().Get(settings.CMaskPsdo),
			rand.Bytes(calcRandSize(psd.fNode.Client().Settings(), size)),
		),
	))
}

func calcRandSize(sett settings.ISettings, len int) uint64 {
	rand := random.NewStdPRNG()
	return uint64(len) + rand.Uint64()%sett.Get(settings.CSizePsdo)
}

func calcRandTime(seconds uint64) time.Duration {
	rand := random.NewStdPRNG()
	return time.Duration(rand.Uint64() % (seconds * 1000))
}
