package network

import (
	"sync"

	"github.com/number571/go-peer/crypto/asymmetric"
)

var (
	_ iF2F = &sF2F{}
)

// F2F connection mode.
type sF2F struct {
	fMutex   sync.Mutex
	fEnabled bool
	fMapping map[string]asymmetric.IPubKey
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
func (f2f *sF2F) InList(pub asymmetric.IPubKey) bool {
	f2f.fMutex.Lock()
	defer f2f.fMutex.Unlock()

	_, ok := f2f.fMapping[pub.Address().String()]
	return ok
}

// Get a list of friends public keys.
func (f2f *sF2F) List() []asymmetric.IPubKey {
	f2f.fMutex.Lock()
	defer f2f.fMutex.Unlock()

	var list []asymmetric.IPubKey
	for _, pub := range f2f.fMapping {
		list = append(list, pub)
	}

	return list
}

// Add public key to list of friends.
func (f2f *sF2F) Append(pub asymmetric.IPubKey) {
	f2f.fMutex.Lock()
	defer f2f.fMutex.Unlock()

	f2f.fMapping[pub.Address().String()] = pub
}

// Delete public key from list of friends.
func (f2f *sF2F) Remove(pub asymmetric.IPubKey) {
	f2f.fMutex.Lock()
	defer f2f.fMutex.Unlock()

	delete(f2f.fMapping, pub.Address().String())
}
