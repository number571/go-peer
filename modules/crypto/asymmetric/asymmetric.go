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

	"github.com/number571/go-peer/modules/crypto/hashing"
)

var (
	_ IPrivKey = &sRSAPrivKey{}
	_ IPubKey  = &sRSAPubKey{}
	_ iAddress = &sAddress{}
)

const (
	cFormatBlock = 32
	CRSAKeyType  = "go-peer/rsa"
)

/*
 * PRIVATE KEY
 */

type sRSAPrivKey struct {
	fPubKey  IPubKey
	fPrivKey *rsa.PrivateKey
}

func newPrivKey(privKey *rsa.PrivateKey) IPrivKey {
	return &sRSAPrivKey{
		fPubKey:  newPubKey(&privKey.PublicKey),
		fPrivKey: privKey,
	}
}

// Create private key by number of bits.
func NewRSAPrivKey(bits uint64) IPrivKey {
	privKey, err := rsa.GenerateKey(rand.Reader, int(bits))
	if err != nil {
		return nil
	}
	return newPrivKey(privKey)
}

func LoadRSAPrivKey(typePrivKey interface{}) IPrivKey {
	switch x := typePrivKey.(type) {
	case []byte:
		privKey := bytesToPrivateKey(x)
		if privKey == nil {
			return nil
		}
		return newPrivKey(privKey)
	case string:
		x = skipSpaceChars(x)
		var (
			prefix = fmt.Sprintf("Priv(%s){", CRSAKeyType)
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
		return LoadRSAPrivKey(pbytes)
	default:
		panic("unsupported type")
	}
}

func (key *sRSAPrivKey) Decrypt(msg []byte) []byte {
	return decryptRSA(key.fPrivKey, msg)
}

func (key *sRSAPrivKey) Sign(msg []byte) []byte {
	return sign(key.fPrivKey, hashing.NewSHA256Hasher(msg).Bytes())
}

func (key *sRSAPrivKey) PubKey() IPubKey {
	return key.fPubKey
}

func (key *sRSAPrivKey) Bytes() []byte {
	return privateKeyToBytes(key.fPrivKey)
}

func (key *sRSAPrivKey) String() string {
	return fmt.Sprintf("Priv(%s){%X}", key.Type(), key.Bytes())
}

func (key *sRSAPrivKey) Format() string {
	s := fmt.Sprintf("Priv(%s){\n", key.Type())
	b := key.Bytes()
	for i := 0; i < len(b); i += cFormatBlock {
		end := i + cFormatBlock
		if end > len(b) {
			end = len(b)
		}
		s += fmt.Sprintf("\t%X\n", b[i:end])
	}
	s += "}"
	return s
}

func (key *sRSAPrivKey) Type() string {
	return CRSAKeyType
}

func (key *sRSAPrivKey) Size() uint64 {
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

type sRSAPubKey struct {
	fAddr   iAddress
	fPubKey *rsa.PublicKey
}

func newPubKey(pubKey *rsa.PublicKey) IPubKey {
	return &sRSAPubKey{
		fAddr:   newAddress(pubKey),
		fPubKey: pubKey,
	}
}

func LoadRSAPubKey(pubkey interface{}) IPubKey {
	switch x := pubkey.(type) {
	case []byte:
		pub := bytesToPublicKey(x)
		if pub == nil {
			return nil
		}
		return newPubKey(pub)
	case string:
		x = skipSpaceChars(x)
		var (
			prefix = fmt.Sprintf("Pub(%s){", CRSAKeyType)
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
		return LoadRSAPubKey(pbytes)
	default:
		panic("unsupported type")
	}
}

func (key *sRSAPubKey) Encrypt(msg []byte) []byte {
	return encryptRSA(key.fPubKey, msg)
}

func (key *sRSAPubKey) Address() iAddress {
	return key.fAddr
}

func (key *sRSAPubKey) Verify(msg []byte, sig []byte) bool {
	return verify(key.fPubKey, hashing.NewSHA256Hasher(msg).Bytes(), sig) == nil
}

func (key *sRSAPubKey) Bytes() []byte {
	return publicKeyToBytes(key.fPubKey)
}

func (key *sRSAPubKey) String() string {
	return fmt.Sprintf("Pub(%s){%X}", key.Type(), key.Bytes())
}

func (key *sRSAPubKey) Format() string {
	s := fmt.Sprintf("Pub(%s){\n", key.Type())
	b := key.Bytes()
	for i := 0; i < len(b); i += cFormatBlock {
		end := i + cFormatBlock
		if end > len(b) {
			end = len(b)
		}
		s += fmt.Sprintf("\t%X\n", b[i:end])
	}
	s += "}"
	return s
}

func (key *sRSAPubKey) Type() string {
	return CRSAKeyType
}

func (key *sRSAPubKey) Size() uint64 {
	return uint64(key.fPubKey.N.BitLen())
}

/*
 * Address
 */

type sAddress struct {
	fBytes []byte
}

func newAddress(pubKey *rsa.PublicKey) iAddress {
	return &sAddress{
		fBytes: hashing.NewSHA256Hasher(
			publicKeyToBytes(pubKey),
		).Bytes(),
	}
}

func (addr *sAddress) Bytes() []byte {
	return addr.fBytes
}

func (addr *sAddress) String() string {
	return fmt.Sprintf("Address(%s){%X}", addr.Type(), addr.Bytes())
}

func (addr *sAddress) Type() string {
	return CRSAKeyType
}

func (addr *sAddress) Size() uint64 {
	return hashing.CSHA256Size
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

func skipSpaceChars(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}