package asymmetric

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/hashing"
)

var (
	_ IPrivKey = &sRSAPrivKey{}
	_ IPubKey  = &sRSAPubKey{}
	_ IAddress = &sAddress{}
)

const (
	cPrivKeyPrefix = "PrivKey{"
	cPubKeyPrefix  = "PubKey{"
	cAddressPrefix = "Address{"
	cKeySuffix     = "}"
)

/*
 * PRIVATE KEY
 */

type sRSAPrivKey struct {
	fPubKey  IPubKey
	fPrivKey *rsa.PrivateKey
}

func newRSAPrivKey(pPrivKey *rsa.PrivateKey) IPrivKey {
	return &sRSAPrivKey{
		fPubKey:  newRSAPubKey(&pPrivKey.PublicKey),
		fPrivKey: pPrivKey,
	}
}

// Create private key by number of bits.
func NewRSAPrivKey(pBits uint64) IPrivKey {
	privKey, err := rsa.GenerateKey(rand.Reader, int(pBits))
	if err != nil {
		panic(err)
	}
	return newRSAPrivKey(privKey)
}

func LoadRSAPrivKey(pPrivKey interface{}) IPrivKey {
	switch x := pPrivKey.(type) {
	case []byte:
		privKey := bytesToPrivateKey(x)
		if privKey == nil {
			return nil
		}
		return newRSAPrivKey(privKey)
	case string:
		x = skipSpaceChars(x)
		var (
			prefix = cPrivKeyPrefix
			suffix = cKeySuffix
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
		return LoadRSAPrivKey(pbytes)
	default:
		panic("unsupported type")
	}
}

func (p *sRSAPrivKey) DecryptBytes(pMsg []byte) []byte {
	return decryptRSA(p.fPrivKey, pMsg)
}

func (p *sRSAPrivKey) SignBytes(pMsg []byte) []byte {
	return sign(p.fPrivKey, hashing.NewSHA256Hasher(pMsg).ToBytes())
}

func (p *sRSAPrivKey) GetPubKey() IPubKey {
	return p.fPubKey
}

func (p *sRSAPrivKey) ToBytes() []byte {
	return privateKeyToBytes(p.fPrivKey)
}

func (p *sRSAPrivKey) ToString() string {
	return fmt.Sprintf(cPrivKeyPrefix+"%X"+cKeySuffix, p.ToBytes())
}

func (p *sRSAPrivKey) GetSize() uint64 {
	return p.GetPubKey().GetSize()
}

// Used PKCS1.
func bytesToPrivateKey(pPrivData []byte) *rsa.PrivateKey {
	priv, err := x509.ParsePKCS1PrivateKey(pPrivData)
	if err != nil {
		return nil
	}
	return priv
}

// Used RSA(OAEP).
func decryptRSA(pPrivKey *rsa.PrivateKey, pData []byte) []byte {
	data, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, pPrivKey, pData, nil)
	if err != nil {
		return nil
	}
	return data
}

// Used PKCS1.
func privateKeyToBytes(pPrivKey *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(pPrivKey)
}

func sign(pPrivKey *rsa.PrivateKey, pHash []byte) []byte {
	signature, err := rsa.SignPSS(rand.Reader, pPrivKey, crypto.SHA256, pHash, nil)
	if err != nil {
		return nil
	}
	return signature
}

/*
 * PUBLIC KEY
 */

type sRSAPubKey struct {
	fAddr   IAddress
	fPubKey *rsa.PublicKey
}

func newRSAPubKey(pPubKey *rsa.PublicKey) IPubKey {
	return &sRSAPubKey{
		fAddr:   newRSAAddress(pPubKey),
		fPubKey: pPubKey,
	}
}

func LoadRSAPubKey(pPubKey interface{}) IPubKey {
	switch x := pPubKey.(type) {
	case []byte:
		pub := bytesToRSAPublicKey(x)
		if pub == nil {
			return nil
		}
		return newRSAPubKey(pub)
	case string:
		x = skipSpaceChars(x)
		var (
			prefix = cPubKeyPrefix
			suffix = cKeySuffix
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
		return LoadRSAPubKey(pbytes)
	default:
		panic("unsupported type")
	}
}

func (p *sRSAPubKey) EncryptBytes(pMsg []byte) []byte {
	return encryptRSA(p.fPubKey, pMsg)
}

func (p *sRSAPubKey) GetAddress() IAddress {
	return p.fAddr
}

func (p *sRSAPubKey) VerifyBytes(pMsg []byte, pSign []byte) bool {
	return verifyRSA(p.fPubKey, hashing.NewSHA256Hasher(pMsg).ToBytes(), pSign) == nil
}

func (p *sRSAPubKey) ToBytes() []byte {
	return rsaPublicKeyToBytes(p.fPubKey)
}

func (p *sRSAPubKey) ToString() string {
	return fmt.Sprintf(cPubKeyPrefix+"%X"+cKeySuffix, p.ToBytes())
}

func (p *sRSAPubKey) GetSize() uint64 {
	return uint64(p.fPubKey.N.BitLen())
}

/*
 * Address
 */

type sAddress struct {
	fBytes []byte
}

func newRSAAddress(pPubKey *rsa.PublicKey) IAddress {
	return &sAddress{
		fBytes: hashing.NewSHA256Hasher(
			rsaPublicKeyToBytes(pPubKey),
		).ToBytes(),
	}
}

func (p *sAddress) ToBytes() []byte {
	return p.fBytes
}

func (p *sAddress) ToString() string {
	return fmt.Sprintf(cAddressPrefix+"%X"+cKeySuffix, p.ToBytes())
}

func (p *sAddress) GetSize() uint64 {
	return hashing.CSHA256Size
}

// Used RSA(OAEP).
func encryptRSA(pPubKey *rsa.PublicKey, pData []byte) []byte {
	data, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pPubKey, pData, nil)
	if err != nil {
		return nil
	}
	return data
}

// Used PKCS1.
func bytesToRSAPublicKey(pPubData []byte) *rsa.PublicKey {
	pub, err := x509.ParsePKCS1PublicKey(pPubData)
	if err != nil {
		return nil
	}
	return pub
}

// Used PKCS1.
func rsaPublicKeyToBytes(pPubKey *rsa.PublicKey) []byte {
	return x509.MarshalPKCS1PublicKey(pPubKey)
}

// Used RSA(PSS).
func verifyRSA(pPubKey *rsa.PublicKey, pHash, pSign []byte) error {
	return rsa.VerifyPSS(pPubKey, crypto.SHA256, pHash, pSign, nil)
}

func skipSpaceChars(pS string) string {
	s := pS
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}
