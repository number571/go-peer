package asymmetric

import (
	"testing"
)

func TestSigner(t *testing.T) {
	t.Parallel()

	privKey := NewDSAPrivKey()
	privKey = LoadDSAPrivKey(privKey.ToBytes())
	if pk := LoadDSAPrivKey([]byte{123}); pk != nil {
		t.Error("success load dsa priv key")
		return
	}

	pubKey := privKey.GetPubKey()
	pubKey = LoadDSAPubKey(pubKey.ToBytes())
	if pk := LoadDSAPubKey([]byte{123}); pk != nil {
		t.Error("success load dsa pub key")
		return
	}

	msg := []byte("hello, world!")
	sign := privKey.SignBytes(msg)

	if !pubKey.VerifyBytes(msg, sign) {
		t.Error("invalid verify")
		return
	}

	// fmt.Println(len(privKey.ToBytes()))
	// fmt.Println(len(pubKey.ToBytes()), len(sign))
}
