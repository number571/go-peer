package routing

import (
	"github.com/number571/go-peer/crypto/asymmetric"
)

var (
	_ IRoute = &sRoute{}
)

// Basic structure for set route to package.
type sRoute struct {
	fReceiver asymmetric.IPubKey
	fPsender  asymmetric.IPrivKey
	fList     []asymmetric.IPubKey
}

// Create route object with receiver.
func NewRoute(receiver asymmetric.IPubKey) IRoute {
	if receiver == nil {
		return nil
	}
	return &sRoute{
		fReceiver: receiver,
	}
}

func (route *sRoute) WithRedirects(psender asymmetric.IPrivKey, list []asymmetric.IPubKey) IRoute {
	return &sRoute{
		fReceiver: route.Receiver(),
		fPsender:  psender,
		fList:     list,
	}
}

// Return receiver's public key.
func (route *sRoute) Receiver() asymmetric.IPubKey {
	return route.fReceiver
}

// Return pseudo sender's private key.
func (route *sRoute) PSender() asymmetric.IPrivKey {
	return route.fPsender
}

// Return all route as list of public keys.
func (route *sRoute) List() []asymmetric.IPubKey {
	return route.fList
}
