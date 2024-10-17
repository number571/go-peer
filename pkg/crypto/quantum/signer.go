package quantum

import (
	"crypto"
	"crypto/rand"

	"github.com/cloudflare/circl/sign"
	dilithium "github.com/cloudflare/circl/sign/dilithium/mode3"
)

const (
	CSignatureSize = dilithium.SignatureSize
)

var (
	_ ISignerPrivKey = &sSignerPrivKey{}
	_ ISignerPubKey  = &sSignerPubKey{}
)

type sSignerPrivKey struct {
	fPrivKey *dilithium.PrivateKey
	fPubKey  *sSignerPubKey
}

type sSignerPubKey struct {
	fPubKey *dilithium.PublicKey
}

func NewSignerPrivKey() ISignerPrivKey {
	_, privKey, err := dilithium.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return newSignerPrivKey(privKey)
}

func LoadSignerPrivKey(pBytes []byte) ISignerPrivKey {
	privKey, err := dilithium.Scheme().UnmarshalBinaryPrivateKey(pBytes)
	if err != nil {
		return nil
	}
	return newSignerPrivKey(privKey)
}

func LoadSignerPubKey(pBytes []byte) ISignerPubKey {
	pubKey, err := dilithium.Scheme().UnmarshalBinaryPublicKey(pBytes)
	if err != nil {
		return nil
	}
	return newSignerPubKey(pubKey)
}

func newSignerPrivKey(privKey sign.PrivateKey) *sSignerPrivKey {
	return &sSignerPrivKey{
		fPrivKey: privKey.(*dilithium.PrivateKey),
		fPubKey:  newSignerPubKey(privKey.Public()),
	}
}

func newSignerPubKey(pubKey crypto.PublicKey) *sSignerPubKey {
	return &sSignerPubKey{pubKey.(*dilithium.PublicKey)}
}

func (p *sSignerPrivKey) GetPubKey() ISignerPubKey {
	return p.fPubKey
}

func (p *sSignerPrivKey) ToBytes() []byte {
	return p.fPrivKey.Bytes()
}

func (p *sSignerPrivKey) SignBytes(pMsg []byte) []byte {
	sign, err := p.fPrivKey.Sign(rand.Reader, pMsg, crypto.Hash(0))
	if err != nil {
		panic(err)
	}
	return sign
}

func (p *sSignerPubKey) VerifyBytes(pMsg, pSign []byte) bool {
	return dilithium.Verify(p.fPubKey, pMsg, pSign)
}

func (p *sSignerPubKey) ToBytes() []byte {
	return p.fPubKey.Bytes()
}
