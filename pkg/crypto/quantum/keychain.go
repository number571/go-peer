package quantum

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
	fKEM         IKEMPrivKey
	fSigner      ISignerPrivKey
	fPubKeyChain IPubKeyChain
}

type sPubKeyChain struct {
	fKEM    IKEMPubKey
	fSigner ISignerPubKey
}

func NewPrivKeyChain(pKEM IKEMPrivKey, pSigner ISignerPrivKey) IPrivKeyChain {
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
	kemPrivKey := LoadKEMPrivKey(pbytesKEM)
	if kemPrivKey == nil {
		return nil
	}

	pbytesSigner, err := hex.DecodeString(splited[1])
	if err != nil {
		return nil
	}
	signerPrivKey := LoadSignerPrivKey(pbytesSigner)
	if signerPrivKey == nil {
		return nil
	}

	return NewPrivKeyChain(kemPrivKey, signerPrivKey)
}

func (p *sPrivKeyChain) ToString() string {
	return fmt.Sprintf("%s%X;%X%s", cPrivKeyPrefix, p.fKEM.ToBytes(), p.fSigner.ToBytes(), cKeySuffix)
}

func (p *sPrivKeyChain) GetKEMPrivKey() IKEMPrivKey {
	return p.fKEM
}

func (p *sPrivKeyChain) GetSignerPrivKey() ISignerPrivKey {
	return p.fSigner
}

func NewPubKeyChain(pKEM IKEMPubKey, pSigner ISignerPubKey) IPubKeyChain {
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
	kemPubKey := LoadKEMPubKey(pbytesKEM)
	if kemPubKey == nil {
		return nil
	}

	pbytesSigner, err := hex.DecodeString(splited[1])
	if err != nil {
		return nil
	}
	signerPubKey := LoadSignerPubKey(pbytesSigner)
	if signerPubKey == nil {
		return nil
	}

	return NewPubKeyChain(kemPubKey, signerPubKey)
}

func (p *sPubKeyChain) ToString() string {
	return fmt.Sprintf("%s%X;%X%s", cPubKeyPrefix, p.fKEM.ToBytes(), p.fSigner.ToBytes(), cKeySuffix)
}

func (p *sPubKeyChain) GetKEMPubKey() IKEMPubKey {
	return p.fKEM
}

func (p *sPubKeyChain) GetSignerPubKey() ISignerPubKey {
	return p.fSigner
}

func skipSpaceChars(pS string) string {
	s := pS
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}
