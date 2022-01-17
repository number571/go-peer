package local

import "github.com/number571/go-peer/crypto"

// Basic structure for set route to package.
type Route struct {
	receiver crypto.PubKey
	psender  crypto.PrivKey
	routes   []crypto.PubKey
}

// Create route object with receiver.
func NewRoute(receiver crypto.PubKey) *Route {
	if receiver == nil {
		return nil
	}
	return &Route{
		receiver: receiver,
	}
}

// Return receiver's public key.
func (route *Route) Receiver() crypto.PubKey {
	return route.receiver
}

// Append pseude sender and routes.
func (route *Route) WithRoad(psender crypto.PrivKey, routes []crypto.PubKey) *Route {
	return &Route{
		receiver: route.receiver,
		psender:  psender,
		routes:   routes,
	}
}
