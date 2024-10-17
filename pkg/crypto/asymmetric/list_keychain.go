package asymmetric

import (
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
)

var (
	_ IListPubKeyChains = &sListPubKeyChains{}
)

// F2F connection mode.
type sListPubKeyChains struct {
	fMutex   sync.RWMutex
	fMapping map[string]IPubKeyChain
}

func NewListPubKeyChains() IListPubKeyChains {
	return &sListPubKeyChains{
		fMapping: make(map[string]IPubKeyChain),
	}
}

// Check the existence of a friend in the list by the public key.
func (p *sListPubKeyChains) GetPubKeyChain(pSignPubKey ISignPubKey) (IPubKeyChain, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	keychain, ok := p.fMapping[hashkey(pSignPubKey)]
	return keychain, ok
}

// Get a list of friends public keys.
func (p *sListPubKeyChains) AllPubKeyChains() []IPubKeyChain {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	list := make([]IPubKeyChain, 0, len(p.fMapping))
	for _, pub := range p.fMapping {
		list = append(list, pub)
	}

	return list
}

// Add public key to list of friends.
func (p *sListPubKeyChains) AddPubKeyChain(pPubKeyChain IPubKeyChain) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fMapping[hashkey(pPubKeyChain.GetSignPubKey())] = pPubKeyChain
}

// Delete public key from list of friends.
func (p *sListPubKeyChains) DelPubKeyChain(pPubKeyChain IPubKeyChain) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	delete(p.fMapping, hashkey(pPubKeyChain.GetSignPubKey()))
}

func hashkey(pSignPubKey ISignPubKey) string {
	return hashing.NewHasher(
		pSignPubKey.ToBytes(),
	).ToString()
}
