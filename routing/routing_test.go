package routing

import (
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
)

func TestRouting(t *testing.T) {
	privKey := asymmetric.NewRSAPrivKey(1024)
	pubKey := privKey.PubKey()

	route := NewRoute(pubKey).
		WithRedirects(privKey, []asymmetric.IPubKey{pubKey})

	if route.Receiver().Address().String() != pubKey.Address().String() {
		t.Error("receiver address not equal address of public key")
		return
	}

	if route.PSender().String() != privKey.String() {
		t.Error("pseudo sender not equal prviate key")
		return
	}

	if route.List()[0].Address().String() != pubKey.Address().String() {
		t.Error("address in list not equal address of public key")
		return
	}
}
