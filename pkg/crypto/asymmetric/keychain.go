package asymmetric

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
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
	fKEM    IKEncPrivKey
	fSigner ISignPrivKey
	fPubKey IPubKey
}

type sPubKey struct {
	fKEM    IKEncPubKey
	fSigner ISignPubKey
	fHasher hashing.IHasher
}

func NewPrivKey() IPrivKey {
	return newPrivKey(NewKEncPrivKey(), NewSignPrivKey())
}

func newPrivKey(pKEM IKEncPrivKey, pSign ISignPrivKey) IPrivKey {
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

	if len(keychainBytes) != (CKEncPrivKeySize + CSignPrivKeySize) {
		return nil
	}

	kemPrivKey := LoadKEncPrivKey(keychainBytes[:CKEncPrivKeySize])
	if kemPrivKey == nil {
		return nil
	}

	signerPrivKey := LoadSignPrivKey(keychainBytes[CKEncPrivKeySize:])
	if signerPrivKey == nil {
		return nil
	}

	return newPrivKey(kemPrivKey, signerPrivKey)
}

func (p *sPrivKey) ToBytes() []byte {
	return bytes.Join([][]byte{p.fKEM.ToBytes(), p.fSigner.ToBytes()}, []byte{})
}

func (p *sPrivKey) ToString() string {
	return fmt.Sprintf("%s%X%s", cPrivKeyPrefix, p.ToBytes(), cKeySuffix)
}

func (p *sPrivKey) GetKEncPrivKey() IKEncPrivKey {
	return p.fKEM
}

func (p *sPrivKey) GetSignPrivKey() ISignPrivKey {
	return p.fSigner
}

func NewPubKey(pKEM IKEncPubKey, pSigner ISignPubKey) IPubKey {
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

	if len(keychainBytes) != (CKEncPubKeySize + CSignPubKeySize) {
		return nil
	}

	kemPubKey := LoadKEncPubKey(keychainBytes[:CKEncPubKeySize])
	if kemPubKey == nil {
		return nil
	}

	signerPubKey := LoadSignPubKey(keychainBytes[CKEncPubKeySize:])
	if signerPubKey == nil {
		return nil
	}

	return NewPubKey(kemPubKey, signerPubKey)
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

func (p *sPubKey) GetKEncPubKey() IKEncPubKey {
	return p.fKEM
}

func (p *sPubKey) GetSignPubKey() ISignPubKey {
	return p.fSigner
}

func skipSpaceChars(pS string) string {
	s := pS
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}
