package asymmetric

import (
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IMapPubKeys = &sMapPubKeys{}
)

// F2F connection mode.
type sMapPubKeys struct {
	fMutex   sync.RWMutex
	fMapping map[string]IPubKey
}

func NewMapPubKeys(pPubKeys ...IPubKey) IMapPubKeys {
	mapPubKeys := &sMapPubKeys{
		fMapping: make(map[string]IPubKey, 256),
	}
	for _, pk := range pPubKeys {
		mapPubKeys.SetPubKey(pk)
	}
	return mapPubKeys
}

// Check the existence of a friend in the list by the public key.
func (p *sMapPubKeys) GetPubKey(pHash []byte) IPubKey {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	pubKey, ok := p.fMapping[encoding.HexEncode(pHash)]
	if !ok {
		return nil
	}
	return pubKey
}

// Delete public key from list of friends.
func (p *sMapPubKeys) DelPubKey(pPubKey IPubKey) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	delete(p.fMapping, pPubKey.GetHasher().ToString())
}

// Add public key to list of friends.
func (p *sMapPubKeys) SetPubKey(pPubKey IPubKey) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fMapping[pPubKey.GetHasher().ToString()] = pPubKey
}
