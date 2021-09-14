package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"

	"github.com/number571/gopeer/encoding"
)

var (
	_ PrivKey = &PrivKeyRSA{}
	_ PubKey  = &PubKeyRSA{}
)

const (
	KeyType       = "gopeer\\rsa"
	TruncatedSize = 20
)

/*
 * PRIVATE KEY
 */

type PrivKeyRSA struct {
	priv *rsa.PrivateKey
}

// Create private key by number of bits.
func NewPrivKey(bits uint) PrivKey {
	priv, err := rsa.GenerateKey(rand.Reader, int(bits))
	if err != nil {
		return nil
	}
	return &PrivKeyRSA{priv}
}

// Create private key by number of bits.
func LoadPrivKey(pbytes []byte) PrivKey {
	return &PrivKeyRSA{bytesToPrivateKey(pbytes)}
}

func (key *PrivKeyRSA) Decrypt(msg []byte) []byte {
	return decryptRSA(key.priv, msg)
}

func (key *PrivKeyRSA) Bytes() []byte {
	return privateKeyToBytes(key.priv)
}

func (key *PrivKeyRSA) String() string {
	return fmt.Sprintf("Priv(%s){%X}", key.Type(), key.Bytes())
}

func (key *PrivKeyRSA) Sign(msg []byte) []byte {
	return sign(key.priv, msg)
}

func (key *PrivKeyRSA) Type() string {
	return KeyType
}

func (key *PrivKeyRSA) PubKey() PubKey {
	return &PubKeyRSA{&key.priv.PublicKey}
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

func sign(priv *rsa.PrivateKey, data []byte) []byte {
	signature, err := rsa.SignPSS(rand.Reader, priv, crypto.SHA256, data, nil)
	if err != nil {
		return nil
	}
	return signature
}

/*
 * PUBLIC KEY
 */

type PubKeyRSA struct {
	pub *rsa.PublicKey
}

func LoadPubKey(pbytes []byte) PubKey {
	return &PubKeyRSA{bytesToPublicKey(pbytes)}
}

func (key *PubKeyRSA) Encrypt(msg []byte) []byte {
	return encryptRSA(key.pub, msg)
}

func (key *PubKeyRSA) Address() Address {
	return Address(hashPublicKey(key.pub))
}

func (key *PubKeyRSA) Bytes() []byte {
	return publicKeyToBytes(key.pub)
}

func (key *PubKeyRSA) String() string {
	return fmt.Sprintf("Pub(%s){%X}", key.Type(), key.Bytes())
}

func (key *PubKeyRSA) Verify(msg []byte, sig []byte) bool {
	return verify(key.pub, msg, sig) == nil
}

func (key *PubKeyRSA) Type() string {
	return KeyType
}

func (key *PubKeyRSA) Size() uint {
	return uint(key.pub.N.BitLen())
}

// Used RSA(OAEP).
func encryptRSA(pub *rsa.PublicKey, data []byte) []byte {
	data, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, data, nil)
	if err != nil {
		return nil
	}
	return data
}

// Hash(PublicKey).
func hashPublicKey(pub *rsa.PublicKey) string {
	return encoding.Base64Encode(HashSum(publicKeyToBytes(pub))[:TruncatedSize])
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
func verify(pub *rsa.PublicKey, data, sign []byte) error {
	return rsa.VerifyPSS(pub, crypto.SHA256, data, sign, nil)
}