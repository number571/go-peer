package asymmetric

import (
	"crypto/rand"

	"github.com/cloudflare/circl/kem"
	mlkem "github.com/cloudflare/circl/kem/mlkem/mlkem768"
)

const (
	CKEMPrivKeySize = mlkem.PrivateKeySize
	CKEMPubKeySize  = mlkem.PublicKeySize
	CKEncSize       = mlkem.CiphertextSize
)

var (
	_ IKEMPrivKey = &sKyber768PrivKey{}
	_ IKEMPubKey  = &sKyber768PubKey{}
)

type sKyber768PrivKey struct {
	fPrivKey *mlkem.PrivateKey
	fPubKey  *sKyber768PubKey
}

type sKyber768PubKey struct {
	fK *mlkem.PublicKey
}

func NewKEMPrivKey() IKEMPrivKey {
	_, privKey, err := mlkem.GenerateKeyPair(rand.Reader)
	if err != nil {
		panic(err)
	}
	return newKEMPrivKey(privKey)
}

func newKEMPrivKey(privKey kem.PrivateKey) *sKyber768PrivKey {
	pk, ok := privKey.(*mlkem.PrivateKey)
	if !ok {
		return nil
	}
	return &sKyber768PrivKey{
		fPrivKey: pk,
		fPubKey:  newKEMPubKey(privKey.Public()),
	}
}

func newKEMPubKey(pubKey kem.PublicKey) *sKyber768PubKey {
	pk, ok := pubKey.(*mlkem.PublicKey)
	if !ok {
		return nil
	}
	return &sKyber768PubKey{pk}
}

func LoadKEMPrivKey(pBytes []byte) IKEMPrivKey {
	privKey, err := mlkem.Scheme().UnmarshalBinaryPrivateKey(pBytes)
	if err != nil {
		return nil
	}
	return newKEMPrivKey(privKey)
}

func LoadKEMPubKey(pBytes []byte) IKEMPubKey {
	pubKey, err := mlkem.Scheme().UnmarshalBinaryPublicKey(pBytes)
	if err != nil {
		return nil
	}
	return newKEMPubKey(pubKey)
}

func (p *sKyber768PubKey) Encapsulate() ([]byte, []byte, error) {
	return mlkem.Scheme().Encapsulate(p.fK)
}

func (p *sKyber768PubKey) ToBytes() []byte {
	b, err := p.fK.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return b
}

func (p *sKyber768PrivKey) GetPubKey() IKEMPubKey {
	return p.fPubKey
}

func (p *sKyber768PrivKey) Decapsulate(pCiphertext []byte) ([]byte, error) {
	return mlkem.Scheme().Decapsulate(p.fPrivKey, pCiphertext)
}

func (p *sKyber768PrivKey) ToBytes() []byte {
	b, err := p.fPrivKey.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return b
}
