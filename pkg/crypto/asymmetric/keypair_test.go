package asymmetric

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/hashing"
)

func TestPanicLoadPrivKey(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = LoadPrivKey(struct{}{})
}

func TestPanicLoadPubKey(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = LoadPubKey(struct{}{})
}

func TestInvalidPrivKey(t *testing.T) {
	t.Parallel()

	if pk := LoadPrivKey("123"); pk != nil {
		t.Error("load priv key (1)")
		return
	}
	if pk := LoadPrivKey(cPrivKeyPrefix); pk != nil {
		t.Error("load priv key (2)")
		return
	}
	if pk := LoadPrivKey(cPrivKeyPrefix + "x" + cKeySuffix); pk != nil {
		t.Error("load priv key (3)")
		return
	}
}

func TestInvalidPubKey(t *testing.T) {
	t.Parallel()

	if pk := LoadPubKey("123"); pk != nil {
		t.Error("load pub key (1)")
		return
	}
	if pk := LoadPubKey(cPubKeyPrefix); pk != nil {
		t.Error("load pub key (2)")
		return
	}
	if pk := LoadPubKey(cPubKeyPrefix + "x" + cKeySuffix); pk != nil {
		t.Error("load pub key (3)")
		return
	}
}

func TestPrivKey(t *testing.T) {
	t.Parallel()

	kemPrivKey := NewKEMPrivKey()
	signerPrivKey := NewDSAPrivKey()

	keychain := newPrivKey(kemPrivKey, signerPrivKey)
	if !bytes.Equal(kemPrivKey.ToBytes(), keychain.GetKEMPrivKey().ToBytes()) {
		t.Error("invalid kem priv key (1)")
		return
	}
	if !bytes.Equal(signerPrivKey.ToBytes(), keychain.GetDSAPrivKey().ToBytes()) {
		t.Error("invalid signer priv key (1)")
		return
	}

	keychain = LoadPrivKey(keychain.ToString())
	if !bytes.Equal(kemPrivKey.ToBytes(), keychain.GetKEMPrivKey().ToBytes()) {
		t.Error("invalid kem priv key (2)")
		return
	}
	if !bytes.Equal(signerPrivKey.ToBytes(), keychain.GetDSAPrivKey().ToBytes()) {
		t.Error("invalid signer priv key (2)")
		return
	}

	keychain = LoadPrivKey(keychain.ToBytes())
	if !bytes.Equal(kemPrivKey.ToBytes(), keychain.GetKEMPrivKey().ToBytes()) {
		t.Error("invalid kem priv key (3)")
		return
	}
	if !bytes.Equal(signerPrivKey.ToBytes(), keychain.GetDSAPrivKey().ToBytes()) {
		t.Error("invalid signer priv key (3)")
		return
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
		return
	}
	if !bytes.Equal(signerPubKey.ToBytes(), keychain.GetDSAPubKey().ToBytes()) {
		t.Error("invalid signer pub key (1)")
		return
	}

	keychain = LoadPubKey(keychain.ToString())
	if !bytes.Equal(kemPubKey.ToBytes(), keychain.GetKEMPubKey().ToBytes()) {
		t.Error("invalid kem pub key (2)")
		return
	}
	if !bytes.Equal(signerPubKey.ToBytes(), keychain.GetDSAPubKey().ToBytes()) {
		t.Error("invalid signer pub key (2)")
		return
	}

	keychain = LoadPubKey(keychain.ToBytes())
	if !bytes.Equal(kemPubKey.ToBytes(), keychain.GetKEMPubKey().ToBytes()) {
		t.Error("invalid kem pub key (3)")
		return
	}
	if !bytes.Equal(signerPubKey.ToBytes(), keychain.GetDSAPubKey().ToBytes()) {
		t.Error("invalid signer pub key (3)")
		return
	}

	if keychain.GetHasher().ToString() != hashing.NewHasher(keychain.ToBytes()).ToString() {
		t.Error("got invalid keychain hasher")
		return
	}

	// fmt.Println("pub.key", len(keychain.ToString()))
}
