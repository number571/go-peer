package asymmetric

import (
	"bytes"
	"fmt"
	"testing"
)

const (
	tcECDHPrivKey = `PrivKey(go-peer/ecdh){3BEA20397555DA5A3FE0BF57C5282A6DD9FE34BE5A0B6605E688852766E24CB1}`
	tcECDHPubKey  = `PubKey(go-peer/ecdh){04ABAD857ADA211613F34658EA58C455DBAD4509FFB7922D34BEA7F7CF245DDF48351A07D09127F63F8AD1FC60F52B2950A904D957C300ED7E4F64C9305BFF43DE}`
)

func TestLoadECDHKey(t *testing.T) {
	privKey := LoadECDHPrivKey(tcECDHPrivKey)
	if privKey == nil {
		t.Error("ecdh privKey key is nil")
		return
	}

	pubKey := LoadECDHPubKey(tcECDHPubKey)
	if pubKey == nil {
		t.Error("ecdh public key is nil")
		return
	}

	if pubKey.ToString() != tcECDHPubKey {
		t.Error("ecdh public key is not equal")
		return
	}

	if privKey.GetPubKey().ToString() != tcECDHPubKey {
		t.Error("ecdh public key from private key is not equal")
		return
	}

	if err := testECDHConverter(privKey, pubKey); err != nil {
		t.Error(err)
		return
	}
}

func testECDHConverter(priv IEphPrivKey, pub IEphPubKey) error {
	if priv.GetSize() != CCurveSize {
		return fmt.Errorf("ecdh private key size != CCurveSize")
	}

	if pub.GetSize() != CCurveSize {
		return fmt.Errorf("ecdh public key size != CCurveSize")
	}

	if priv.ToString() != tcECDHPrivKey {
		return fmt.Errorf("ecdh private key string != tcECDHPrivKey")
	}

	if pub.ToString() != tcECDHPubKey {
		return fmt.Errorf("ecdh public key string != tcECDHPubKey")
	}

	return nil
}

func TestECDHSharedKey(t *testing.T) {
	var (
		priv1 = NewECDHPrivKey()
		pub1  = priv1.GetPubKey()

		priv2 = LoadECDHPrivKey(tcECDHPrivKey)
		pub2  = priv2.GetPubKey()
	)

	sharedKey1, err := priv1.GetSharedKey(pub2)
	if err != nil {
		t.Error(err)
		return
	}

	sharedKey2, err := priv2.GetSharedKey(pub1)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(sharedKey1, sharedKey2) {
		t.Error("got invalid shared keys")
		return
	}
}
