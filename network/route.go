package network

import (
	"github.com/number571/gopeer/crypto"
)

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

// Append pseudo sender to route.
func (route *Route) WithSender(psender crypto.PrivKey) *Route {
	route.psender = psender
	return route
}

// Append route-nodes to route.
func (route *Route) WithRoutes(routes []crypto.PubKey) *Route {
	route.routes = routes
	return route
}
