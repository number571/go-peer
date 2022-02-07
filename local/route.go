package local

import "github.com/number571/go-peer/crypto"

// Basic structure for set route to package.
type RouteT struct {
	receiver crypto.PubKey
	psender  crypto.PrivKey
	routes   []crypto.PubKey
}

// Create route object with receiver.
func NewRoute(receiver crypto.PubKey, psender crypto.PrivKey, routes []crypto.PubKey) *RouteT {
	if receiver == nil {
		return nil
	}
	return &RouteT{
		receiver: receiver,
		psender:  psender,
		routes:   routes,
	}
}

// Return receiver's public key.
func (route *RouteT) Receiver() crypto.PubKey {
	return route.receiver
}
