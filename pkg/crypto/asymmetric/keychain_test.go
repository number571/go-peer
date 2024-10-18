package asymmetric

import (
	"bytes"
	"testing"
)

func TestPrivKey(t *testing.T) {
	t.Parallel()

	kemPrivKey := NewKEncPrivKey()
	signerPrivKey := NewSignPrivKey()

	keychain := newPrivKey(kemPrivKey, signerPrivKey)
	if !bytes.Equal(kemPrivKey.ToBytes(), keychain.GetKEncPrivKey().ToBytes()) {
		t.Error("invalid kem priv key (1)")
	}
	if !bytes.Equal(signerPrivKey.ToBytes(), keychain.GetSignPrivKey().ToBytes()) {
		t.Error("invalid signer priv key (1)")
	}

	keychain = LoadPrivKey(keychain.ToString())
	if !bytes.Equal(kemPrivKey.ToBytes(), keychain.GetKEncPrivKey().ToBytes()) {
		t.Error("invalid kem priv key (2)")
	}
	if !bytes.Equal(signerPrivKey.ToBytes(), keychain.GetSignPrivKey().ToBytes()) {
		t.Error("invalid signer priv key (2)")
	}

	// fmt.Println("priv.key", len(keychain.ToString()))
}

func TestPubKey(t *testing.T) {
	t.Parallel()

	kemPubKey := NewKEncPrivKey().GetPubKey()
	signerPubKey := NewSignPrivKey().GetPubKey()

	keychain := NewPubKey(kemPubKey, signerPubKey)
	if !bytes.Equal(kemPubKey.ToBytes(), keychain.GetKEncPubKey().ToBytes()) {
		t.Error("invalid kem pub key (1)")
	}
	if !bytes.Equal(signerPubKey.ToBytes(), keychain.GetSignPubKey().ToBytes()) {
		t.Error("invalid signer pub key (1)")
	}

	keychain = LoadPubKey(keychain.ToString())
	if !bytes.Equal(kemPubKey.ToBytes(), keychain.GetKEncPubKey().ToBytes()) {
		t.Error("invalid kem pub key (2)")
	}
	if !bytes.Equal(signerPubKey.ToBytes(), keychain.GetSignPubKey().ToBytes()) {
		t.Error("invalid signer pub key (2)")
	}

	// fmt.Println("pub.key", len(keychain.ToString()))
}
