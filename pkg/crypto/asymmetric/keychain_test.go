package asymmetric

import (
	"bytes"
	"testing"
)

func TestPrivKey(t *testing.T) {
	t.Parallel()

	kemPrivKey := NewKEMPrivKey()
	signerPrivKey := NewDSAPrivKey()

	keychain := newPrivKey(kemPrivKey, signerPrivKey)
	if !bytes.Equal(kemPrivKey.ToBytes(), keychain.GetKEMPrivKey().ToBytes()) {
		t.Error("invalid kem priv key (1)")
	}
	if !bytes.Equal(signerPrivKey.ToBytes(), keychain.GetDSAPrivKey().ToBytes()) {
		t.Error("invalid signer priv key (1)")
	}

	keychain = LoadPrivKey(keychain.ToString())
	if !bytes.Equal(kemPrivKey.ToBytes(), keychain.GetKEMPrivKey().ToBytes()) {
		t.Error("invalid kem priv key (2)")
	}
	if !bytes.Equal(signerPrivKey.ToBytes(), keychain.GetDSAPrivKey().ToBytes()) {
		t.Error("invalid signer priv key (2)")
	}

	// fmt.Println("priv.key", len(keychain.ToString()))
}

func TestPubKey(t *testing.T) {
	t.Parallel()

	kemPubKey := NewKEMPrivKey().GetPubKey()
	signerPubKey := NewDSAPrivKey().GetPubKey()

	keychain := NewPubKey(kemPubKey, signerPubKey)
	if !bytes.Equal(kemPubKey.ToBytes(), keychain.GetKEMPubKey().ToBytes()) {
		t.Error("invalid kem pub key (1)")
	}
	if !bytes.Equal(signerPubKey.ToBytes(), keychain.GetDSAPubKey().ToBytes()) {
		t.Error("invalid signer pub key (1)")
	}

	keychain = LoadPubKey(keychain.ToString())
	if !bytes.Equal(kemPubKey.ToBytes(), keychain.GetKEMPubKey().ToBytes()) {
		t.Error("invalid kem pub key (2)")
	}
	if !bytes.Equal(signerPubKey.ToBytes(), keychain.GetDSAPubKey().ToBytes()) {
		t.Error("invalid signer pub key (2)")
	}

	// fmt.Println("pub.key", len(keychain.ToString()))
}
