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
func (l *sListPubKeys) InPubKeys(pub IPubKey) bool {
	l.fMutex.Lock()
	defer l.fMutex.Unlock()

	_, ok := l.fMapping[pub.Address().ToString()]
	return ok
}

// Get a list of friends public keys.
func (l *sListPubKeys) GetPubKeys() []IPubKey {
	l.fMutex.Lock()
	defer l.fMutex.Unlock()

	var list []IPubKey
	for _, pub := range l.fMapping {
		list = append(list, pub)
	}

	return list
}

// Add public key to list of friends.
func (l *sListPubKeys) AddPubKey(pub IPubKey) {
	l.fMutex.Lock()
	defer l.fMutex.Unlock()

	l.fMapping[pub.Address().ToString()] = pub
}

// Delete public key from list of friends.
func (l *sListPubKeys) DelPubKey(pub IPubKey) {
	l.fMutex.Lock()
	defer l.fMutex.Unlock()

	delete(l.fMapping, pub.Address().ToString())
}
