package network

import (
	"bytes"
	"sync"
	"time"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/settings"
)

var (
	_ iChecker     = &sChecker{}
	_ iCheckerInfo = &sCheckerInfo{}
)

type sChecker struct {
	fMutex   sync.Mutex
	fNode    INode
	fEnabled bool
	fChannel chan struct{}
	fMapping map[string]iCheckerInfo
}

type sCheckerInfo struct {
	fMutex  sync.Mutex
	fOnline bool
	fPubKey crypto.IPubKey
}

// Set state = bool.
func (checker *sChecker) Switch(state bool) {
	checker.fMutex.Lock()
	defer checker.fMutex.Unlock()

	if checker.fEnabled == state {
		return
	}
	checker.fEnabled = state

	switch checker.fEnabled {
	case true:
		checker.start()
	case false:
		checker.stop()
	}
}

// Get current state of online mode.
func (checker *sChecker) Status() bool {
	checker.fMutex.Lock()
	defer checker.fMutex.Unlock()

	return checker.fEnabled
}

// Check the existence in the list by the public key.
func (checker *sChecker) InList(pub crypto.IPubKey) bool {
	checker.fMutex.Lock()
	defer checker.fMutex.Unlock()

	_, ok := checker.fMapping[pub.Address()]
	return ok
}

// Get a list of checks public keys.
func (checker *sChecker) List() []crypto.IPubKey {
	checker.fMutex.Lock()
	defer checker.fMutex.Unlock()

	var list []crypto.IPubKey
	for _, chk := range checker.fMapping {
		list = append(list, chk.PubKey())
	}

	return list
}

// Get a list of checks public keys and online status.
func (checker *sChecker) ListWithInfo() []iCheckerInfo {
	checker.fMutex.Lock()
	defer checker.fMutex.Unlock()

	var list []iCheckerInfo
	for _, chk := range checker.fMapping {
		list = append(list, chk)
	}

	return list
}

// Add public key to list of checks.
func (checker *sChecker) Append(pub crypto.IPubKey) {
	checker.fMutex.Lock()
	defer checker.fMutex.Unlock()

	checker.fMapping[pub.Address()] = &sCheckerInfo{
		fOnline: false,
		fPubKey: pub,
	}
}

// Delete public key from list of checks.
func (checker *sChecker) Remove(pub crypto.IPubKey) {
	checker.fMutex.Lock()
	defer checker.fMutex.Unlock()

	delete(checker.fMapping, pub.Address())
}

func (checkerInfo *sCheckerInfo) PubKey() crypto.IPubKey {
	checkerInfo.fMutex.Lock()
	defer checkerInfo.fMutex.Unlock()

	return checkerInfo.fPubKey
}

func (checkerInfo *sCheckerInfo) Online() bool {
	checkerInfo.fMutex.Lock()
	defer checkerInfo.fMutex.Unlock()

	return checkerInfo.fOnline
}

func (checker *sChecker) start() {
	sett := checker.fNode.Client().Settings()
	patt := encoding.Uint64ToBytes(sett.Get(settings.MaskPing))
	go func(checker *sChecker, patt []byte, timeOut, timePing uint64) {
		node := checker.fNode.(*sNode)
		for {
			for _, recv := range checker.ListWithInfo() {
				go func(recv *sCheckerInfo) {
					resp, err := node.doRequest(
						local.NewRoute(recv.fPubKey),
						local.NewMessage(patt, patt),
						1, // retry number
						timeOut,
					)
					if err != nil || !bytes.Equal(resp, patt) {
						recv.fOnline = false
						return
					}
					recv.fOnline = true
				}(recv.(*sCheckerInfo))
			}
			select {
			case <-checker.fChannel:
				return
			case <-time.After(time.Second * time.Duration(timePing)):
				continue
			}
		}
	}(checker, patt, sett.Get(settings.TimeWait), sett.Get(settings.TimePing))
}

func (checker *sChecker) stop() {
	checker.fChannel <- struct{}{}
}
