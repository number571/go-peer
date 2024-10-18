package asymmetric

import (
	"bytes"
	"testing"

	"github.com/cloudflare/circl/kem"
)

func TestNewKEM(t *testing.T) {
	t.Parallel()

	if pk := newKEMPrivKey(&tsPrivateKeyKEM{}); pk != nil {
		t.Error("success get another kem privkey (not mlkem768)")
		return
	}

	if pk := newKEMPubKey(&tsPublicKeyKEM{}); pk != nil {
		t.Error("success get another kem pubkey (not mlkem768)")
		return
	}
}

func TestKEM(t *testing.T) {
	t.Parallel()

	privKey := NewKEMPrivKey()
	privKey = LoadKEMPrivKey(privKey.ToBytes())
	if pk := LoadKEMPrivKey([]byte{123}); pk != nil {
		t.Error("success load kem priv key")
		return
	}

	pubKey := privKey.GetPubKey()
	pubKey = LoadKEMPubKey(pubKey.ToBytes())
	if pk := LoadKEMPubKey([]byte{123}); pk != nil {
		t.Error("success load kem pub key")
		return
	}

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

	if !bytes.Equal(ss1, ss2) {
		t.Error("invalid shared secret")
		return
	}

	// fmt.Println(len(privKey.ToBytes()))
	// fmt.Println(len(pubKey.ToBytes()), len(ct), len(ss1))
}

var (
	_ kem.PrivateKey = &tsPrivateKeyKEM{}
	_ kem.PublicKey  = &tsPublicKeyKEM{}
)

type tsPrivateKeyKEM struct{}
type tsPublicKeyKEM struct{}

func (p *tsPrivateKeyKEM) Scheme() kem.Scheme             { return nil }
func (p *tsPrivateKeyKEM) MarshalBinary() ([]byte, error) { return nil, nil }
func (p *tsPrivateKeyKEM) Equal(kem.PrivateKey) bool      { return false }
func (p *tsPrivateKeyKEM) Public() kem.PublicKey          { return nil }

func (p *tsPublicKeyKEM) Scheme() kem.Scheme             { return nil }
func (p *tsPublicKeyKEM) MarshalBinary() ([]byte, error) { return nil, nil }
func (p *tsPublicKeyKEM) Equal(kem.PublicKey) bool       { return false }
