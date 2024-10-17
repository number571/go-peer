package quantum

import (
	"crypto/rand"

	"github.com/cloudflare/circl/kem"
	kyber "github.com/cloudflare/circl/kem/kyber/kyber768"
)

const (
	CCiphertextSize = kyber.CiphertextSize
)

var (
	_ IKEMPrivKey = &sKEMPrivKey{}
	_ IKEMPubKey  = &sKEMPubKey{}
)

type sKEMPrivKey struct {
	fPrivKey *kyber.PrivateKey
	fPubKey  *sKEMPubKey
}

type sKEMPubKey struct {
	fK *kyber.PublicKey
}

func NewKEMPrivKey() IKEMPrivKey {
	_, privKey, err := kyber.GenerateKeyPair(rand.Reader)
	if err != nil {
		panic(err)
	}
	return newKEMPrivKey(privKey)
}

func newKEMPrivKey(privKey kem.PrivateKey) *sKEMPrivKey {
	return &sKEMPrivKey{
		fPrivKey: privKey.(*kyber.PrivateKey),
		fPubKey:  newKEMPubKey(privKey.Public()),
	}
}

func newKEMPubKey(pubKey kem.PublicKey) *sKEMPubKey {
	return &sKEMPubKey{pubKey.(*kyber.PublicKey)}
}

func LoadKEMPrivKey(pBytes []byte) IKEMPrivKey {
	privKey, err := kyber.Scheme().UnmarshalBinaryPrivateKey(pBytes)
	if err != nil {
		return nil
	}
	return newKEMPrivKey(privKey)
}

func LoadKEMPubKey(pBytes []byte) IKEMPubKey {
	pubKey, err := kyber.Scheme().UnmarshalBinaryPublicKey(pBytes)
	if err != nil {
		return nil
	}
	return newKEMPubKey(pubKey)
}

func (p *sKEMPubKey) Encapsulate() ([]byte, []byte, error) {
	return kyber.Scheme().Encapsulate(p.fK)
}

func (p *sKEMPubKey) ToBytes() []byte {
	b, err := p.fK.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return b
}

func (p *sKEMPrivKey) GetPubKey() IKEMPubKey {
	return p.fPubKey
}

func (p *sKEMPrivKey) Decapsulate(pCiphertext []byte) ([]byte, error) {
	return kyber.Scheme().Decapsulate(p.fPrivKey, pCiphertext)
}

func (p *sKEMPrivKey) ToBytes() []byte {
	b, err := p.fPrivKey.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return b
}
