package asymmetric

import (
	"crypto/rand"

	"github.com/cloudflare/circl/kem"
	kyber "github.com/cloudflare/circl/kem/kyber/kyber768"
)

const (
	CKEncPrivKeySize = kyber.PrivateKeySize
	CKEncPubKeySize  = kyber.PublicKeySize
	CKEncSize        = kyber.CiphertextSize
)

var (
	_ IKEncPrivKey = &sKyber768PrivKey{}
	_ IKEncPubKey  = &sKyber768PubKey{}
)

type sKyber768PrivKey struct {
	fPrivKey *kyber.PrivateKey
	fPubKey  *sKyber768PubKey
}

type sKyber768PubKey struct {
	fK *kyber.PublicKey
}

func NewKEncPrivKey() IKEncPrivKey {
	_, privKey, err := kyber.GenerateKeyPair(rand.Reader)
	if err != nil {
		panic(err)
	}
	return newKEncPrivKey(privKey)
}

func newKEncPrivKey(privKey kem.PrivateKey) *sKyber768PrivKey {
	pk, ok := privKey.(*kyber.PrivateKey)
	if !ok {
		return nil
	}
	return &sKyber768PrivKey{
		fPrivKey: pk,
		fPubKey:  newKEncPubKey(privKey.Public()),
	}
}

func newKEncPubKey(pubKey kem.PublicKey) *sKyber768PubKey {
	pk, ok := pubKey.(*kyber.PublicKey)
	if !ok {
		return nil
	}
	return &sKyber768PubKey{pk}
}

func LoadKEncPrivKey(pBytes []byte) IKEncPrivKey {
	privKey, err := kyber.Scheme().UnmarshalBinaryPrivateKey(pBytes)
	if err != nil {
		return nil
	}
	return newKEncPrivKey(privKey)
}

func LoadKEncPubKey(pBytes []byte) IKEncPubKey {
	pubKey, err := kyber.Scheme().UnmarshalBinaryPublicKey(pBytes)
	if err != nil {
		return nil
	}
	return newKEncPubKey(pubKey)
}

func (p *sKyber768PubKey) Encapsulate() ([]byte, []byte, error) {
	return kyber.Scheme().Encapsulate(p.fK)
}

func (p *sKyber768PubKey) ToBytes() []byte {
	b, err := p.fK.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return b
}

func (p *sKyber768PrivKey) GetPubKey() IKEncPubKey {
	return p.fPubKey
}

func (p *sKyber768PrivKey) Decapsulate(pCiphertext []byte) ([]byte, error) {
	return kyber.Scheme().Decapsulate(p.fPrivKey, pCiphertext)
}

func (p *sKyber768PrivKey) ToBytes() []byte {
	b, err := p.fPrivKey.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return b
}
