package asymmetric

import (
	"testing"
)

func TestMapPubKeys(t *testing.T) {
	t.Parallel()

	pubKeys := []IPubKey{
		NewPrivKey().GetPubKey(),
		NewPrivKey().GetPubKey(),
		NewPrivKey().GetPubKey(),
	}

	mapping := NewMapPubKeys(pubKeys...)
	pkHash := pubKeys[1].GetHasher().ToString()
	if pk := mapping.GetPubKey(pkHash); pk == nil {
		t.Error("get invalid pub key")
		return
	}

	mapping.DelPubKey(pubKeys[1])
	if pk := mapping.GetPubKey(pkHash); pk != nil {
		t.Error("get success deleted pub key")
		return
	}
}
