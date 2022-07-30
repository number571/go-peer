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

func (psd *sPseudo) request() iPseudo {
	psd.fMutex.Lock()
	defer psd.fMutex.Unlock()

	if !psd.fEnabled {
		return psd
	}

	psd.doRequest()
	return psd
}

func (psd *sPseudo) sleep() iPseudo {
	psd.fMutex.Lock()
	defer psd.fMutex.Unlock()

	if !psd.fEnabled {
		return psd
	}

	sett := psd.fNode.Client().Settings()
	rslp := sett.Get(settings.CTimeRslp)
	time.Sleep(calcRandTimeInMS(0, rslp))
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
	wait := sett.Get(settings.CTimePreq)
	go func(psd *sPseudo, wait uint64) {
		for {
			psd.doRequest()
			select {
			case <-psd.fChannel:
				return
			case <-time.After(calcRandTimeInMS(1, wait)):
				continue
			}
		}
	}(psd, wait)
}

func (psd *sPseudo) stop() {
	psd.fChannel <- struct{}{}
}

func (psd *sPseudo) doRequest() {
	sett := psd.fNode.Client().Settings()
	rand := random.NewStdPRNG()
	psd.fNode.Broadcast(psd.fNode.Client().Encrypt(
		routing.NewRoute(psd.fPrivKey.PubKey()),
		payload.NewPayload(
			sett.Get(settings.CMaskPsdo),
			rand.Bytes(calcRandSize(sett)),
		),
	))
}

func calcRandSize(sett settings.ISettings) uint64 {
	rand := random.NewStdPRNG()
	sizePack := sett.Get(settings.CSizePack)
	return rand.Uint64() % (sizePack / 2)
}

func calcRandTimeInMS(minSeconds, addMaxSeconds uint64) time.Duration {
	mnum := minSeconds * 1000
	rnum := random.NewStdPRNG().Uint64() % (addMaxSeconds * 1000)
	return time.Millisecond * time.Duration(int64(mnum)+int64(rnum))
}
