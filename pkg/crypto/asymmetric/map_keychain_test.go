package asymmetric

import (
	"bytes"
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
		list.SetPubKey(pk.GetDSAPubKey(), pk.GetKEMPubKey())
	}

	dsaPubKey := pubKeys[1].GetDSAPubKey()
	pk, ok := list.GetPubKey(dsaPubKey)
	if !ok || !bytes.Equal(pk.ToBytes(), pubKeys[1].GetKEMPubKey().ToBytes()) {
		t.Error("get invalid pub key")
		return
	}

	list.DelPubKey(dsaPubKey)
	if _, ok := list.GetPubKey(dsaPubKey); ok {
		t.Error("get success deleted pub key")
		return
	}
}
