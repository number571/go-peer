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
	_ IPrivKeyChain = &sPrivKeyChain{}
	_ IPubKeyChain  = &sPubKeyChain{}
)

type sPrivKeyChain struct {
	fKEM         IKEncPrivKey
	fSigner      ISignPrivKey
	fPubKeyChain IPubKeyChain
}

type sPubKeyChain struct {
	fKEM    IKEncPubKey
	fSigner ISignPubKey
	fHasher hashing.IHasher
}

func NewPrivKeyChain(pKEM IKEncPrivKey, pSigner ISignPrivKey) IPrivKeyChain {
	return &sPrivKeyChain{
		fKEM:         pKEM,
		fSigner:      pSigner,
		fPubKeyChain: NewPubKeyChain(pKEM.GetPubKey(), pSigner.GetPubKey()),
	}
}

func (p *sPrivKeyChain) GetPubKeyChain() IPubKeyChain {
	return p.fPubKeyChain
}

func LoadPrivKeyChain(pKeychain interface{}) IPrivKeyChain {
	pKeychainBytes := []byte{}

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
		pKeychainBytes = encoding.HexDecode(s)
	case []byte:
		pKeychainBytes = x
	default:
		panic("unknown type private key chain")
	}

	if len(pKeychainBytes) != (CKEncPrivKeySize + CSignPrivKeySize) {
		return nil
	}

	kemPrivKey := LoadKEncPrivKey(pKeychainBytes[:CKEncPrivKeySize])
	if kemPrivKey == nil {
		return nil
	}

	signerPrivKey := LoadSignPrivKey(pKeychainBytes[CKEncPrivKeySize:])
	if signerPrivKey == nil {
		return nil
	}

	return NewPrivKeyChain(kemPrivKey, signerPrivKey)
}

func (p *sPrivKeyChain) ToBytes() []byte {
	return bytes.Join([][]byte{p.fKEM.ToBytes(), p.fSigner.ToBytes()}, []byte{})
}

func (p *sPrivKeyChain) ToString() string {
	return fmt.Sprintf("%s%X%s", cPrivKeyPrefix, p.ToBytes(), cKeySuffix)
}

func (p *sPrivKeyChain) GetKEncPrivKey() IKEncPrivKey {
	return p.fKEM
}

func (p *sPrivKeyChain) GetSignPrivKey() ISignPrivKey {
	return p.fSigner
}

func NewPubKeyChain(pKEM IKEncPubKey, pSigner ISignPubKey) IPubKeyChain {
	pubKeyChain := &sPubKeyChain{
		fKEM:    pKEM,
		fSigner: pSigner,
	}
	pubKeyChain.fHasher = hashing.NewHasher(pubKeyChain.ToBytes())
	return pubKeyChain
}

func LoadPubKeyChain(pKeychain interface{}) IPubKeyChain {
	pKeychainBytes := []byte{}

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
		pKeychainBytes = encoding.HexDecode(s)
	case []byte:
		pKeychainBytes = x
	default:
		panic("unknown type public key chain")
	}

	if len(pKeychainBytes) != (CKEncPubKeySize + CSignPubKeySize) {
		return nil
	}

	kemPubKey := LoadKEncPubKey(pKeychainBytes[:CKEncPubKeySize])
	if kemPubKey == nil {
		return nil
	}

	signerPubKey := LoadSignPubKey(pKeychainBytes[CKEncPubKeySize:])
	if signerPubKey == nil {
		return nil
	}

	return NewPubKeyChain(kemPubKey, signerPubKey)
}

func (p *sPubKeyChain) GetHasher() hashing.IHasher {
	return p.fHasher
}

func (p *sPubKeyChain) ToBytes() []byte {
	return bytes.Join([][]byte{p.fKEM.ToBytes(), p.fSigner.ToBytes()}, []byte{})
}

func (p *sPubKeyChain) ToString() string {
	return fmt.Sprintf("%s%X%s", cPubKeyPrefix, p.ToBytes(), cKeySuffix)
}

func (p *sPubKeyChain) GetKEncPubKey() IKEncPubKey {
	return p.fKEM
}

func (p *sPubKeyChain) GetSignPubKey() ISignPubKey {
	return p.fSigner
}

func skipSpaceChars(pS string) string {
	s := pS
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}
