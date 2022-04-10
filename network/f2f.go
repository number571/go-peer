package network

import (
	"sync"

	"github.com/number571/go-peer/crypto"
)

var (
	_ iF2F = &sF2F{}
)

// F2F connection mode.
type sF2F struct {
	fMutex   sync.Mutex
	fEnabled bool
	fMapping map[string]crypto.IPubKey
}

// Set state = bool.
func (f2f *sF2F) Switch(state bool) {
	f2f.fMutex.Lock()
	defer f2f.fMutex.Unlock()

	f2f.fEnabled = state
}

// Get current state of f2f mode.
func (f2f *sF2F) Status() bool {
	f2f.fMutex.Lock()
	defer f2f.fMutex.Unlock()

	return f2f.fEnabled
}

// Check the existence of a friend in the list by the public key.
func (f2f *sF2F) InList(pub crypto.IPubKey) bool {
	f2f.fMutex.Lock()
	defer f2f.fMutex.Unlock()

	_, ok := f2f.fMapping[pub.Address()]
	return ok
}

// Get a list of friends public keys.
func (f2f *sF2F) List() []crypto.IPubKey {
	f2f.fMutex.Lock()
	defer f2f.fMutex.Unlock()

	var list []crypto.IPubKey
	for _, pub := range f2f.fMapping {
		list = append(list, pub)
	}

	return list
}

// Add public key to list of friends.
func (f2f *sF2F) Append(pub crypto.IPubKey) {
	f2f.fMutex.Lock()
	defer f2f.fMutex.Unlock()

	f2f.fMapping[pub.Address()] = pub
}

// Delete public key from list of friends.
func (f2f *sF2F) Remove(pub crypto.IPubKey) {
	f2f.fMutex.Lock()
	defer f2f.fMutex.Unlock()

	delete(f2f.fMapping, pub.Address())
}
