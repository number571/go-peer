package asymmetric

import (
	"sync"
)

var (
	_ IListPubKeys = &sListPubKeys{}
)

// F2F connection mode.
type sListPubKeys struct {
	fMutex   sync.Mutex
	fMapping map[string]IPubKey
}

func NewListPubKeys() IListPubKeys {
	return &sListPubKeys{
		fMapping: make(map[string]IPubKey),
	}
}

// Check the existence of a friend in the list by the public key.
func (p *sListPubKeys) InPubKeys(pPubKey IPubKey) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	_, ok := p.fMapping[pPubKey.GetAddress().ToString()]
	return ok
}

// Get a list of friends public keys.
func (p *sListPubKeys) GetPubKeys() []IPubKey {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	var list []IPubKey
	for _, pub := range p.fMapping {
		list = append(list, pub)
	}

	return list
}

// Add public key to list of friends.
func (p *sListPubKeys) AddPubKey(pPubKey IPubKey) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fMapping[pPubKey.GetAddress().ToString()] = pPubKey
}

// Delete public key from list of friends.
func (p *sListPubKeys) DelPubKey(pub IPubKey) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	delete(p.fMapping, pub.GetAddress().ToString())
}
