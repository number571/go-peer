package asymmetric

import (
	"crypto"
	"io"
	"testing"

	"github.com/cloudflare/circl/sign"
)

func TestNewDSA(t *testing.T) {
	t.Parallel()

	if pk := newDSAPrivKey(&tsPrivateKeyDSA{}); pk != nil {
		t.Error("success get another dsa privkey (not mldsa65)")
		return
	}

	if pk := newDSAPubKey(&tsPublicKeyDSA{}); pk != nil {
		t.Error("success get another dsa pubkey (not mldsa65)")
		return
	}
}

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

var (
	_ sign.PrivateKey = &tsPrivateKeyDSA{}
	_ sign.PublicKey  = &tsPublicKeyDSA{}
)

type tsPrivateKeyDSA struct{}
type tsPublicKeyDSA struct{}

func (p *tsPrivateKeyDSA) Scheme() sign.Scheme            { return nil }
func (p *tsPrivateKeyDSA) MarshalBinary() ([]byte, error) { return nil, nil }
func (p *tsPrivateKeyDSA) Equal(crypto.PrivateKey) bool   { return false }
func (p *tsPrivateKeyDSA) Public() crypto.PublicKey       { return nil }
func (p *tsPrivateKeyDSA) Sign(io.Reader, []byte, crypto.SignerOpts) ([]byte, error) {
	return nil, nil
}

func (p *tsPublicKeyDSA) Scheme() sign.Scheme            { return nil }
func (p *tsPublicKeyDSA) MarshalBinary() ([]byte, error) { return nil, nil }
func (p *tsPublicKeyDSA) Equal(crypto.PublicKey) bool    { return false }
