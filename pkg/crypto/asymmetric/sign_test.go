package asymmetric

import (
	"testing"
)

func TestSigner(t *testing.T) {
	t.Parallel()

	privKey := NewDSAPrivKey()
	privKey = LoadDSAPrivKey(privKey.ToBytes())

	pubKey := privKey.GetPubKey()
	pubKey = LoadDSAPubKey(pubKey.ToBytes())

	msg := []byte("hello, world!")
	sign := privKey.SignBytes(msg)

	if !pubKey.VerifyBytes(msg, sign) {
		t.Error("invalid verify")
		return
	}

	// fmt.Println(len(privKey.ToBytes()))
	// fmt.Println(len(pubKey.ToBytes()), len(sign))
}
