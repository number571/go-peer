package network

import (
	"sync"

	"github.com/number571/go-peer/crypto"
)

var (
	_ F2F = &F2FT{}
)

// F2F connection mode.
type F2FT struct {
	mutex   sync.Mutex
	enabled bool
	friends map[string]crypto.PubKey
}

// Get current state of f2f mode.
func (f2f *F2FT) Status() bool {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()

	return f2f.enabled
}

// Switch f2f mode to reverse.
func (f2f *F2FT) Switch() {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()

	f2f.enabled = !f2f.enabled
}

// Check the existence of a friend in the list by the public key.
func (f2f *F2FT) InList(pub crypto.PubKey) bool {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()

	_, ok := f2f.friends[string(pub.Address())]
	return ok
}

// Get a list of friends public keys.
func (f2f *F2FT) List() []crypto.PubKey {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()

	var list []crypto.PubKey
	for _, pub := range f2f.friends {
		list = append(list, pub)
	}

	return list
}

// Add public key to list of friends.
func (f2f *F2FT) Append(pub crypto.PubKey) {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()

	f2f.friends[string(pub.Address())] = pub
}

// Delete public key from list of friends.
func (f2f *F2FT) Remove(pub crypto.PubKey) {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()

	delete(f2f.friends, string(pub.Address()))
}
