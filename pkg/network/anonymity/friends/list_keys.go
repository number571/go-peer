package friends

import (
	"sync"
)

var (
	_ IListKeys = &sListKeys{}
)

// F2F connection mode.
type sListKeys struct {
	fMutex   sync.RWMutex
	fMapping map[string]struct{}
}

func NewListKeys() IListKeys {
	return &sListKeys{
		fMapping: make(map[string]struct{}, 256),
	}
}

// Get a list of friends public keys.
func (p *sListKeys) GetKeys() [][]byte {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	list := make([][]byte, 0, len(p.fMapping))
	for key := range p.fMapping {
		list = append(list, []byte(key))
	}

	return list
}

// Add public key to list of friends.
func (p *sListKeys) AddKey(pKey []byte) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fMapping[string(pKey)] = struct{}{}
}

// Delete public key from list of friends.
func (p *sListKeys) DelKey(pKey []byte) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	delete(p.fMapping, string(pKey))
}
