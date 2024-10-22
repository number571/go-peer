package asymmetric

import (
	"crypto"
	"crypto/rand"

	"github.com/cloudflare/circl/sign"
	mldsa "github.com/cloudflare/circl/sign/mldsa/mldsa65"
)

const (
	CDSAPrivKeySize = mldsa.PrivateKeySize
	CDSAPubKeySize  = mldsa.PublicKeySize
	CDSAKeySeedSize = mldsa.SeedSize
	CDSASignSize    = mldsa.SignatureSize
)

var (
	_ IDSAPrivKey = &sDilithiumM3PrivKey{}
	_ IDSAPubKey  = &sDilithiumM3PubKey{}
)

type sDilithiumM3PrivKey struct {
	fPrivKey *mldsa.PrivateKey
	fPubKey  *sDilithiumM3PubKey
}

type sDilithiumM3PubKey struct {
	fK *mldsa.PublicKey
}

func NewDSAPrivKeyFromSeed(pSeed []byte) IDSAPrivKey {
	if len(pSeed) != CDSAKeySeedSize {
		panic("len(pSeed) != CDSAKeySeedSize")
	}
	arr := &[mldsa.SeedSize]byte{}
	copy(arr[:], pSeed)
	_, privKey := mldsa.NewKeyFromSeed(arr)
	return newDSAPrivKey(privKey)
}

func NewDSAPrivKey() IDSAPrivKey {
	_, privKey, err := mldsa.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return newDSAPrivKey(privKey)
}

func LoadDSAPrivKey(pBytes []byte) IDSAPrivKey {
	privKey, err := mldsa.Scheme().UnmarshalBinaryPrivateKey(pBytes)
	if err != nil {
		return nil
	}
	return newDSAPrivKey(privKey)
}

func LoadDSAPubKey(pBytes []byte) IDSAPubKey {
	pubKey, err := mldsa.Scheme().UnmarshalBinaryPublicKey(pBytes)
	if err != nil {
		return nil
	}
	return newDSAPubKey(pubKey)
}

func newDSAPrivKey(privKey sign.PrivateKey) *sDilithiumM3PrivKey {
	pk, ok := privKey.(*mldsa.PrivateKey)
	if !ok {
		return nil
	}
	return &sDilithiumM3PrivKey{
		fPrivKey: pk,
		fPubKey:  newDSAPubKey(privKey.Public()),
	}
}

func newDSAPubKey(pubKey crypto.PublicKey) *sDilithiumM3PubKey {
	pk, ok := pubKey.(*mldsa.PublicKey)
	if !ok {
		return nil
	}
	return &sDilithiumM3PubKey{pk}
}

func (p *sDilithiumM3PrivKey) GetPubKey() IDSAPubKey {
	return p.fPubKey
}

func (p *sDilithiumM3PrivKey) ToBytes() []byte {
	return p.fPrivKey.Bytes()
}

func (p *sDilithiumM3PrivKey) SignBytes(pMsg []byte) []byte {
	sign, _ := p.fPrivKey.Sign(rand.Reader, pMsg, crypto.Hash(0))
	return sign
}

func (p *sDilithiumM3PubKey) VerifyBytes(pMsg, pSign []byte) bool {
	return mldsa.Verify(p.fK, pMsg, nil, pSign)
}

func (p *sDilithiumM3PubKey) ToBytes() []byte {
	return p.fK.Bytes()
}
