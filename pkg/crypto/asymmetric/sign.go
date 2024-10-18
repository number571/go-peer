package asymmetric

import (
	"crypto"
	"crypto/rand"

	"github.com/cloudflare/circl/sign"
	dilithium "github.com/cloudflare/circl/sign/dilithium/mode3"
)

const (
	CSignPrivKeySize = dilithium.PrivateKeySize
	CSignPubKeySize  = dilithium.PublicKeySize
	CSignSize        = dilithium.SignatureSize
)

var (
	_ ISignPrivKey = &sDilithiumM3PrivKey{}
	_ ISignPubKey  = &sDilithiumM3PubKey{}
)

type sDilithiumM3PrivKey struct {
	fPrivKey *dilithium.PrivateKey
	fPubKey  *sDilithiumM3PubKey
}

type sDilithiumM3PubKey struct {
	fK *dilithium.PublicKey
}

func NewSignPrivKey() ISignPrivKey {
	_, privKey, err := dilithium.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return newSignPrivKey(privKey)
}

func LoadSignPrivKey(pBytes []byte) ISignPrivKey {
	privKey, err := dilithium.Scheme().UnmarshalBinaryPrivateKey(pBytes)
	if err != nil {
		return nil
	}
	return newSignPrivKey(privKey)
}

func LoadSignPubKey(pBytes []byte) ISignPubKey {
	pubKey, err := dilithium.Scheme().UnmarshalBinaryPublicKey(pBytes)
	if err != nil {
		return nil
	}
	return newSignPubKey(pubKey)
}

func newSignPrivKey(privKey sign.PrivateKey) *sDilithiumM3PrivKey {
	pk, ok := privKey.(*dilithium.PrivateKey)
	if !ok {
		return nil
	}
	return &sDilithiumM3PrivKey{
		fPrivKey: pk,
		fPubKey:  newSignPubKey(privKey.Public()),
	}
}

func newSignPubKey(pubKey crypto.PublicKey) *sDilithiumM3PubKey {
	pk, ok := pubKey.(*dilithium.PublicKey)
	if !ok {
		return nil
	}
	return &sDilithiumM3PubKey{pk}
}

func (p *sDilithiumM3PrivKey) GetPubKey() ISignPubKey {
	return p.fPubKey
}

func (p *sDilithiumM3PrivKey) ToBytes() []byte {
	return p.fPrivKey.Bytes()
}

func (p *sDilithiumM3PrivKey) SignBytes(pMsg []byte) []byte {
	sign, err := p.fPrivKey.Sign(rand.Reader, pMsg, crypto.Hash(0))
	if err != nil {
		panic(err)
	}
	return sign
}

func (p *sDilithiumM3PubKey) VerifyBytes(pMsg, pSign []byte) bool {
	return dilithium.Verify(p.fK, pMsg, pSign)
}

func (p *sDilithiumM3PubKey) ToBytes() []byte {
	return p.fK.Bytes()
}
