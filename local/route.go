package local

import "github.com/number571/go-peer/crypto"

var (
	_ IRoute = &sRoute{}
)

// Basic structure for set route to package.
type sRoute struct {
	fReceiver crypto.IPubKey
	fPsender  crypto.IPrivKey
	fList     []crypto.IPubKey
}

// Create route object with receiver.
func NewRoute(receiver crypto.IPubKey) IRoute {
	if receiver == nil {
		return nil
	}
	return &sRoute{
		fReceiver: receiver,
	}
}

func (route *sRoute) WithRedirects(psender crypto.IPrivKey, list []crypto.IPubKey) IRoute {
	return &sRoute{
		fReceiver: route.Receiver(),
		fPsender:  psender,
		fList:     list,
	}
}

// Return receiver's public key.
func (route *sRoute) Receiver() crypto.IPubKey {
	return route.fReceiver
}

// Return pseudo sender's private key.
func (route *sRoute) PSender() crypto.IPrivKey {
	return route.fPsender
}

// Return all route as list of public keys.
func (route *sRoute) List() []crypto.IPubKey {
	return route.fList
}
