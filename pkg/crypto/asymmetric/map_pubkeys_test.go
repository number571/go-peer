package asymmetric

import (
	"testing"
)

func TestMapPubKeys(t *testing.T) {
	t.Parallel()

	list := NewMapPubKeys()

	pubKeys := []IPubKey{
		NewPrivKey().GetPubKey(),
		NewPrivKey().GetPubKey(),
		NewPrivKey().GetPubKey(),
	}

	for _, pk := range pubKeys {
		list.SetPubKey(pk)
	}

	if ok := list.InPubKeys(pubKeys[1]); !ok {
		t.Error("get invalid pub key")
		return
	}

	list.DelPubKey(pubKeys[1])
	if ok := list.InPubKeys(pubKeys[1]); ok {
		t.Error("get success deleted pub key")
		return
	}
}
