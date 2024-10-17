package asymmetric

import (
	"bytes"
	"testing"
)

func TestPrivKeyChain(t *testing.T) {
	t.Parallel()

	kemPrivKey := NewKEncPrivKey()
	signerPrivKey := NewSignPrivKey()

	keychain := NewPrivKeyChain(kemPrivKey, signerPrivKey)
	if !bytes.Equal(kemPrivKey.ToBytes(), keychain.GetKEncPrivKey().ToBytes()) {
		t.Error("invalid kem priv key (1)")
	}
	if !bytes.Equal(signerPrivKey.ToBytes(), keychain.GetSignPrivKey().ToBytes()) {
		t.Error("invalid signer priv key (1)")
	}

	keychain = LoadPrivKeyChain(keychain.ToString())
	if !bytes.Equal(kemPrivKey.ToBytes(), keychain.GetKEncPrivKey().ToBytes()) {
		t.Error("invalid kem priv key (2)")
	}
	if !bytes.Equal(signerPrivKey.ToBytes(), keychain.GetSignPrivKey().ToBytes()) {
		t.Error("invalid signer priv key (2)")
	}

	// fmt.Println(keychain.ToString())
}

func TestPubKeyChain(t *testing.T) {
	t.Parallel()

	kemPubKey := NewKEncPrivKey().GetPubKey()
	signerPubKey := NewSignPrivKey().GetPubKey()

	keychain := NewPubKeyChain(kemPubKey, signerPubKey)
	if !bytes.Equal(kemPubKey.ToBytes(), keychain.GetKEncPubKey().ToBytes()) {
		t.Error("invalid kem pub key (1)")
	}
	if !bytes.Equal(signerPubKey.ToBytes(), keychain.GetSignPubKey().ToBytes()) {
		t.Error("invalid signer pub key (1)")
	}

	keychain = LoadPubKeyChain(keychain.ToString())
	if !bytes.Equal(kemPubKey.ToBytes(), keychain.GetKEncPubKey().ToBytes()) {
		t.Error("invalid kem pub key (2)")
	}
	if !bytes.Equal(signerPubKey.ToBytes(), keychain.GetSignPubKey().ToBytes()) {
		t.Error("invalid signer pub key (2)")
	}

	// fmt.Println(keychain.ToString())
}
