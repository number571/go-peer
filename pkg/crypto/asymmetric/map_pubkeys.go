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
	fMapping map[string]IKEMPubKey
}

func NewMapPubKeys() IMapPubKeys {
	return &sMapPubKeys{
		fMapping: make(map[string]IKEMPubKey),
	}
}

// Add public key to list of friends.
func (p *sMapPubKeys) SetPubKey(pDSAPubKey IDSAPubKey, pKEMPubKey IKEMPubKey) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fMapping[hashkey(pDSAPubKey)] = pKEMPubKey
}

// Check the existence of a friend in the list by the public key.
func (p *sMapPubKeys) GetPubKey(pDSAPubKey IDSAPubKey) (IKEMPubKey, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	kemPubKey, ok := p.fMapping[hashkey(pDSAPubKey)]
	return kemPubKey, ok
}

// Delete public key from list of friends.
func (p *sMapPubKeys) DelPubKey(pPubKey IDSAPubKey) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	delete(p.fMapping, hashkey(pPubKey))
}

func hashkey(pDSAPubKey IDSAPubKey) string {
	return hashing.NewHasher(pDSAPubKey.ToBytes()).ToString()
}
