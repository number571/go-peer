package local

import "github.com/number571/go-peer/crypto"

// Basic structure for set route to package.
type routeT struct {
	receiver crypto.PubKey
	psender  crypto.PrivKey
	routes   []crypto.PubKey
}

// Create route object with receiver.
func NewRoute(receiver crypto.PubKey, psender crypto.PrivKey, routes []crypto.PubKey) *routeT {
	if receiver == nil {
		return nil
	}
	return &routeT{
		receiver: receiver,
		psender:  psender,
		routes:   routes,
	}
}

// Return receiver's public key.
func (route *routeT) Receiver() crypto.PubKey {
	return route.receiver
}
