package routing

import (
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
)

func TestDefault(t *testing.T) {
	privKey := asymmetric.NewRSAPrivKey(1024)
	pubKey := privKey.PubKey()

	route := NewRoute(pubKey).
		WithRedirects(privKey, []asymmetric.IPubKey{pubKey})

	if route.Receiver().Address() != pubKey.Address() {
		t.Error("receiver address not equal address of public key")
		return
	}

	if route.PSender().String() != privKey.String() {
		t.Error("pseudo sender not equal prviate key")
		return
	}

	if route.List()[0].Address() != pubKey.Address() {
		t.Error("address in list not equal address of public key")
		return
	}
}
