package asymmetric

import (
	"encoding/hex"
	"fmt"
	"strings"
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

func LoadPrivKeyChain(pKeychain string) IPrivKeyChain {
	x := skipSpaceChars(pKeychain)
	if !strings.HasPrefix(x, cPrivKeyPrefix) {
		return nil
	}
	x = strings.TrimPrefix(x, cPrivKeyPrefix)

	if !strings.HasSuffix(x, cKeySuffix) {
		return nil
	}
	x = strings.TrimSuffix(x, cKeySuffix)

	splited := strings.Split(x, ";")
	if len(splited) != 2 {
		return nil
	}

	pbytesKEM, err := hex.DecodeString(splited[0])
	if err != nil {
		return nil
	}
	kemPrivKey := LoadKEncPrivKey(pbytesKEM)
	if kemPrivKey == nil {
		return nil
	}

	pbytesSigner, err := hex.DecodeString(splited[1])
	if err != nil {
		return nil
	}
	signerPrivKey := LoadSignPrivKey(pbytesSigner)
	if signerPrivKey == nil {
		return nil
	}

	return NewPrivKeyChain(kemPrivKey, signerPrivKey)
}

func (p *sPrivKeyChain) ToString() string {
	return fmt.Sprintf("%s%X;%X%s", cPrivKeyPrefix, p.fKEM.ToBytes(), p.fSigner.ToBytes(), cKeySuffix)
}

func (p *sPrivKeyChain) GetKEncPrivKey() IKEncPrivKey {
	return p.fKEM
}

func (p *sPrivKeyChain) GetSignPrivKey() ISignPrivKey {
	return p.fSigner
}

func NewPubKeyChain(pKEM IKEncPubKey, pSigner ISignPubKey) IPubKeyChain {
	return &sPubKeyChain{
		fKEM:    pKEM,
		fSigner: pSigner,
	}
}

func LoadPubKeyChain(pKeychain string) IPubKeyChain {
	x := skipSpaceChars(pKeychain)
	if !strings.HasPrefix(x, cPubKeyPrefix) {
		return nil
	}
	x = strings.TrimPrefix(x, cPubKeyPrefix)

	if !strings.HasSuffix(x, cKeySuffix) {
		return nil
	}
	x = strings.TrimSuffix(x, cKeySuffix)

	splited := strings.Split(x, ";")
	if len(splited) != 2 {
		return nil
	}

	pbytesKEM, err := hex.DecodeString(splited[0])
	if err != nil {
		return nil
	}
	kemPubKey := LoadKEncPubKey(pbytesKEM)
	if kemPubKey == nil {
		return nil
	}

	pbytesSigner, err := hex.DecodeString(splited[1])
	if err != nil {
		return nil
	}
	signerPubKey := LoadSignPubKey(pbytesSigner)
	if signerPubKey == nil {
		return nil
	}

	return NewPubKeyChain(kemPubKey, signerPubKey)
}

func (p *sPubKeyChain) ToString() string {
	return fmt.Sprintf("%s%X;%X%s", cPubKeyPrefix, p.fKEM.ToBytes(), p.fSigner.ToBytes(), cKeySuffix)
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
