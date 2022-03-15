package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"strings"
)

var (
	_ IPrivKey = &sPrivKey{}
	_ IPubKey  = &sPubKey{}
)

const (
	AsymmKeyType = "go-peer\\rsa"
)

/*
 * PRIVATE KEY
 */

type sPrivKey struct {
	priv *rsa.PrivateKey
}

// Create private key by number of bits.
func NewPrivKey(bits uint64) IPrivKey {
	priv, err := rsa.GenerateKey(rand.Reader, int(bits))
	if err != nil {
		return nil
	}
	return &sPrivKey{priv}
}

func LoadPrivKey(privkey interface{}) IPrivKey {
	switch x := privkey.(type) {
	case []byte:
		priv := bytesToPrivateKey(x)
		if priv == nil {
			return nil
		}
		return &sPrivKey{priv}
	case string:
		var (
			prefix = fmt.Sprintf("Priv(%s){", AsymmKeyType)
			suffix = "}"
		)

		if !strings.HasPrefix(x, prefix) {
			return nil
		}
		x = strings.TrimPrefix(x, prefix)

		if !strings.HasSuffix(x, suffix) {
			return nil
		}
		x = strings.TrimSuffix(x, suffix)

		pbytes, err := hex.DecodeString(x)
		if err != nil {
			return nil
		}
		return LoadPrivKey(pbytes)
	default:
		panic("unsupported type")
	}
}

func (key *sPrivKey) Decrypt(msg []byte) []byte {
	return decryptRSA(key.priv, msg)
}

func (key *sPrivKey) Sign(msg []byte) []byte {
	return sign(key.priv, NewHasher(msg).Bytes())
}

func (key *sPrivKey) PubKey() IPubKey {
	return &sPubKey{&key.priv.PublicKey}
}

func (key *sPrivKey) Bytes() []byte {
	return privateKeyToBytes(key.priv)
}

func (key *sPrivKey) String() string {
	return fmt.Sprintf("Priv(%s){%X}", AsymmKeyType, key.Bytes())
}

func (key *sPrivKey) Type() string {
	return AsymmKeyType
}

func (key *sPrivKey) Size() uint64 {
	return key.PubKey().Size()
}

// Used PKCS1.
func bytesToPrivateKey(privData []byte) *rsa.PrivateKey {
	priv, err := x509.ParsePKCS1PrivateKey(privData)
	if err != nil {
		return nil
	}
	return priv
}

// Used RSA(OAEP).
func decryptRSA(priv *rsa.PrivateKey, data []byte) []byte {
	data, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, data, nil)
	if err != nil {
		return nil
	}
	return data
}

// Used PKCS1.
func privateKeyToBytes(priv *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(priv)
}

func sign(priv *rsa.PrivateKey, hash []byte) []byte {
	signature, err := rsa.SignPSS(rand.Reader, priv, crypto.SHA256, hash, nil)
	if err != nil {
		return nil
	}
	return signature
}

/*
 * PUBLIC KEY
 */

type sPubKey struct {
	pub *rsa.PublicKey
}

func LoadPubKey(pubkey interface{}) IPubKey {
	switch x := pubkey.(type) {
	case []byte:
		pub := bytesToPublicKey(x)
		if pub == nil {
			return nil
		}
		return &sPubKey{pub}
	case string:
		var (
			prefix = fmt.Sprintf("Pub(%s){", AsymmKeyType)
			suffix = "}"
		)

		if !strings.HasPrefix(x, prefix) {
			return nil
		}
		x = strings.TrimPrefix(x, prefix)

		if !strings.HasSuffix(x, suffix) {
			return nil
		}
		x = strings.TrimSuffix(x, suffix)

		pbytes, err := hex.DecodeString(x)
		if err != nil {
			return nil
		}
		return LoadPubKey(pbytes)
	default:
		panic("unsupported type")
	}
}

func (key *sPubKey) Encrypt(msg []byte) []byte {
	return encryptRSA(key.pub, msg)
}

func (key *sPubKey) Address() string {
	return NewHasher(key.Bytes()).String()
}

func (key *sPubKey) Verify(msg []byte, sig []byte) bool {
	return verify(key.pub, NewHasher(msg).Bytes(), sig) == nil
}

func (key *sPubKey) Bytes() []byte {
	return publicKeyToBytes(key.pub)
}

func (key *sPubKey) String() string {
	return fmt.Sprintf("Pub(%s){%X}", AsymmKeyType, key.Bytes())
}

func (key *sPubKey) Type() string {
	return AsymmKeyType
}

func (key *sPubKey) Size() uint64 {
	return uint64(key.pub.N.BitLen())
}

// Used RSA(OAEP).
func encryptRSA(pub *rsa.PublicKey, data []byte) []byte {
	data, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, data, nil)
	if err != nil {
		return nil
	}
	return data
}

// Used PKCS1.
func bytesToPublicKey(pubData []byte) *rsa.PublicKey {
	pub, err := x509.ParsePKCS1PublicKey(pubData)
	if err != nil {
		return nil
	}
	return pub
}

// Used PKCS1.
func publicKeyToBytes(pub *rsa.PublicKey) []byte {
	return x509.MarshalPKCS1PublicKey(pub)
}

// Used RSA(PSS).
func verify(pub *rsa.PublicKey, hash, sign []byte) error {
	return rsa.VerifyPSS(pub, crypto.SHA256, hash, sign, nil)
}
