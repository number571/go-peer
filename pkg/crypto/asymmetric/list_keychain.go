package asymmetric

import (
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
)

var (
	_ IListPubKeys = &sListPubKeys{}
)

// F2F connection mode.
type sListPubKeys struct {
	fMutex   sync.RWMutex
	fMapping map[string]IPubKey
}

func NewListPubKeys() IListPubKeys {
	return &sListPubKeys{
		fMapping: make(map[string]IPubKey),
	}
}

// Check the existence of a friend in the list by the public key.
func (p *sListPubKeys) GetPubKey(pSignPubKey ISignPubKey) (IPubKey, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	keychain, ok := p.fMapping[hashkey(pSignPubKey)]
	return keychain, ok
}

// Get a list of friends public keys.
func (p *sListPubKeys) AllPubKeys() []IPubKey {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	list := make([]IPubKey, 0, len(p.fMapping))
	for _, pub := range p.fMapping {
		list = append(list, pub)
	}

	return list
}

// Add public key to list of friends.
func (p *sListPubKeys) AddPubKey(pPubKey IPubKey) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fMapping[hashkey(pPubKey.GetSignPubKey())] = pPubKey
}

// Delete public key from list of friends.
func (p *sListPubKeys) DelPubKey(pPubKey IPubKey) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	delete(p.fMapping, hashkey(pPubKey.GetSignPubKey()))
}

func hashkey(pSignPubKey ISignPubKey) string {
	return hashing.NewHasher(
		pSignPubKey.ToBytes(),
	).ToString()
}
