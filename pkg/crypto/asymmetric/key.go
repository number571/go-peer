package asymmetric

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
)

const (
	CKeySeedSize = CKEMKeySeedSize + CDSAKeySeedSize
	CPrivKeySize = CKEMPrivKeySize + CDSAPrivKeySize
	CPubKeySize  = CKEMPubKeySize + CDSAPubKeySize
)

const (
	cPrivKeyPrefix = "PrivKey{"
	cPubKeyPrefix  = "PubKey{"
	cKeySuffix     = "}"
)

var (
	_ IPrivKey = &sPrivKey{}
	_ IPubKey  = &sPubKey{}
)

type sPrivKey struct {
	fKEM    IKEMPrivKey
	fSigner IDSAPrivKey
	fPubKey IPubKey
}

type sPubKey struct {
	fKEM    IKEMPubKey
	fSigner IDSAPubKey
	fHasher hashing.IHasher
}

func NewPrivKeyFromSeed(pSeed []byte) IPrivKey {
	if len(pSeed) != CKeySeedSize {
		panic("len(pSeed) != CKeySeedSize")
	}
	return newPrivKey(
		NewKEMPrivKeyFromSeed(pSeed[:CKEMKeySeedSize]),
		NewDSAPrivKeyFromSeed(pSeed[CKEMKeySeedSize:]),
	)
}

func NewPrivKey() IPrivKey {
	return newPrivKey(NewKEMPrivKey(), NewDSAPrivKey())
}

func newPrivKey(pKEM IKEMPrivKey, pSign IDSAPrivKey) IPrivKey {
	return &sPrivKey{
		fKEM:    pKEM,
		fSigner: pSign,
		fPubKey: NewPubKey(pKEM.GetPubKey(), pSign.GetPubKey()),
	}
}

func (p *sPrivKey) GetPubKey() IPubKey {
	return p.fPubKey
}

func LoadPrivKey(pKeychain interface{}) IPrivKey {
	keychainBytes := []byte{} // nolint: staticcheck

	switch x := pKeychain.(type) {
	case string:
		s := skipSpaceChars(x)
		if !strings.HasPrefix(s, cPrivKeyPrefix) {
			return nil
		}
		s = strings.TrimPrefix(s, cPrivKeyPrefix)
		if !strings.HasSuffix(s, cKeySuffix) {
			return nil
		}
		s = strings.TrimSuffix(s, cKeySuffix)
		keychainBytes = encoding.HexDecode(s)
	case []byte:
		keychainBytes = x
	default:
		panic("unknown type private key chain")
	}

	if len(keychainBytes) != CPrivKeySize {
		return nil
	}

	return newPrivKey(
		LoadKEMPrivKey(keychainBytes[:CKEMPrivKeySize]),
		LoadDSAPrivKey(keychainBytes[CKEMPrivKeySize:]),
	)
}

func (p *sPrivKey) ToBytes() []byte {
	return bytes.Join([][]byte{p.fKEM.ToBytes(), p.fSigner.ToBytes()}, []byte{})
}

func (p *sPrivKey) ToString() string {
	return fmt.Sprintf("%s%X%s", cPrivKeyPrefix, p.ToBytes(), cKeySuffix)
}

func (p *sPrivKey) GetKEMPrivKey() IKEMPrivKey {
	return p.fKEM
}

func (p *sPrivKey) GetDSAPrivKey() IDSAPrivKey {
	return p.fSigner
}

func NewPubKey(pKEM IKEMPubKey, pSigner IDSAPubKey) IPubKey {
	pubKeyChain := &sPubKey{
		fKEM:    pKEM,
		fSigner: pSigner,
	}
	pubKeyChain.fHasher = hashing.NewHasher(pubKeyChain.ToBytes())
	return pubKeyChain
}

func LoadPubKey(pKeychain interface{}) IPubKey {
	keychainBytes := []byte{} // nolint: staticcheck

	switch x := pKeychain.(type) {
	case string:
		s := skipSpaceChars(x)
		if !strings.HasPrefix(s, cPubKeyPrefix) {
			return nil
		}
		s = strings.TrimPrefix(s, cPubKeyPrefix)
		if !strings.HasSuffix(s, cKeySuffix) {
			return nil
		}
		s = strings.TrimSuffix(s, cKeySuffix)
		keychainBytes = encoding.HexDecode(s)
	case []byte:
		keychainBytes = x
	default:
		panic("unknown type public key chain")
	}

	if len(keychainBytes) != CPubKeySize {
		return nil
	}

	return NewPubKey(
		LoadKEMPubKey(keychainBytes[:CKEMPubKeySize]),
		LoadDSAPubKey(keychainBytes[CKEMPubKeySize:]),
	)
}

func (p *sPubKey) GetHasher() hashing.IHasher {
	return p.fHasher
}

func (p *sPubKey) ToBytes() []byte {
	return bytes.Join([][]byte{p.fKEM.ToBytes(), p.fSigner.ToBytes()}, []byte{})
}

func (p *sPubKey) ToString() string {
	return fmt.Sprintf("%s%X%s", cPubKeyPrefix, p.ToBytes(), cKeySuffix)
}

func (p *sPubKey) GetKEMPubKey() IKEMPubKey {
	return p.fKEM
}

func (p *sPubKey) GetDSAPubKey() IDSAPubKey {
	return p.fSigner
}

func skipSpaceChars(pS string) string {
	s := pS
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}
