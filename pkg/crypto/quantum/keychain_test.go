package quantum

import (
	"bytes"
	"testing"
)

func TestPrivKeyChain(t *testing.T) {
	t.Parallel()

	kemPrivKey := NewKEMPrivKey()
	signerPrivKey := NewSignerPrivKey()

	keychain := NewPrivKeyChain(kemPrivKey, signerPrivKey)
	if !bytes.Equal(kemPrivKey.ToBytes(), keychain.GetKEMPrivKey().ToBytes()) {
		t.Error("invalid kem priv key (1)")
	}
	if !bytes.Equal(signerPrivKey.ToBytes(), keychain.GetSignerPrivKey().ToBytes()) {
		t.Error("invalid signer priv key (1)")
	}

	keychain = LoadPrivKeyChain(keychain.ToString())
	if !bytes.Equal(kemPrivKey.ToBytes(), keychain.GetKEMPrivKey().ToBytes()) {
		t.Error("invalid kem priv key (2)")
	}
	if !bytes.Equal(signerPrivKey.ToBytes(), keychain.GetSignerPrivKey().ToBytes()) {
		t.Error("invalid signer priv key (2)")
	}

	// fmt.Println(keychain.ToString())
}

func TestPubKeyChain(t *testing.T) {
	t.Parallel()

	kemPubKey := NewKEMPrivKey().GetPubKey()
	signerPubKey := NewSignerPrivKey().GetPubKey()

	keychain := NewPubKeyChain(kemPubKey, signerPubKey)
	if !bytes.Equal(kemPubKey.ToBytes(), keychain.GetKEMPubKey().ToBytes()) {
		t.Error("invalid kem pub key (1)")
	}
	if !bytes.Equal(signerPubKey.ToBytes(), keychain.GetSignerPubKey().ToBytes()) {
		t.Error("invalid signer pub key (1)")
	}

	keychain = LoadPubKeyChain(keychain.ToString())
	if !bytes.Equal(kemPubKey.ToBytes(), keychain.GetKEMPubKey().ToBytes()) {
		t.Error("invalid kem pub key (2)")
	}
	if !bytes.Equal(signerPubKey.ToBytes(), keychain.GetSignerPubKey().ToBytes()) {
		t.Error("invalid signer pub key (2)")
	}

	// fmt.Println(keychain.ToString())
}
