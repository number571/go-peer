package asymmetric

import (
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
)

var (
	_ IMapPubKeys = &sMapPubKeys{}
)

// F2F connection mode.
type sMapPubKeys struct {
	fMutex   sync.RWMutex
	fMapping map[string]struct{}
}

func NewMapPubKeys() IMapPubKeys {
	return &sMapPubKeys{
		fMapping: make(map[string]struct{}),
	}
}

// Add public key to list of friends.
func (p *sMapPubKeys) SetPubKey(pPubKey IPubKey) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fMapping[hashkey(pPubKey)] = struct{}{}
}

// Check the existence of a friend in the list by the public key.
func (p *sMapPubKeys) InPubKeys(pPubKey IPubKey) bool {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	_, ok := p.fMapping[hashkey(pPubKey)]
	return ok
}

// Delete public key from list of friends.
func (p *sMapPubKeys) DelPubKey(pPubKey IPubKey) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	delete(p.fMapping, hashkey(pPubKey))
}

func hashkey(pPubKey IPubKey) string {
	return hashing.NewHasher(pPubKey.ToBytes()).ToString()
}
