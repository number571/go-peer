package quantum

import (
	"bytes"
	"testing"
)

func TestKEM(t *testing.T) {
	t.Parallel()

	privKey := NewKEMPrivKey()
	privKey = LoadKEMPrivKey(privKey.ToBytes())

	pubKey := privKey.GetPubKey()
	pubKey = LoadKEMPubKey(pubKey.ToBytes())

	ct, ss1, err := pubKey.Encapsulate()
	if err != nil {
		t.Error(err)
		return
	}

	ss2, err := privKey.Decapsulate(ct)
	if err != nil {
		t.Error(err)
		return
	}

	// fmt.Println(len(pubKey.ToBytes()), len(ct), len(ss1))

	if !bytes.Equal(ss1, ss2) {
		t.Error("invalid shared secret")
		return
	}
}
