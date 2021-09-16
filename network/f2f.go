package network

import (
	"sync"

	"github.com/number571/gopeer/crypto"
)

// F2F connection mode.
type friendToFriend struct {
	mutex   sync.Mutex
	enabled bool
	friends map[string]crypto.PubKey
}

// Get current state of f2f mode.
func (f2f *friendToFriend) State() bool {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	return f2f.enabled
}

// Switch f2f mode to reverse.
func (f2f *friendToFriend) Switch() {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	f2f.enabled = !f2f.enabled
}

// Check the existence of a friend in the list by the public key.
func (f2f *friendToFriend) InList(pub crypto.PubKey) bool {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	_, ok := f2f.friends[string(pub.Address())]
	return ok
}

// Get a list of friends public keys.
func (f2f *friendToFriend) List() []crypto.PubKey {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	var list []crypto.PubKey
	for _, pub := range f2f.friends {
		list = append(list, pub)
	}
	return list
}

// Add public key to list of friends.
func (f2f *friendToFriend) Append(pubs ...crypto.PubKey) {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	for _, pub := range pubs {
		f2f.friends[string(pub.Address())] = pub
	}
}

// Delete public key from list of friends.
func (f2f *friendToFriend) Remove(pubs ...crypto.PubKey) {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	for _, pub := range pubs {
		delete(f2f.friends, string(pub.Address()))
	}
}
