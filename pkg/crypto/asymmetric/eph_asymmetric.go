package asymmetric

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/number571/go-peer/pkg/errors"
)

var (
	_ IEphPrivKey = &sECDHPrivKey{}
	_ IEphPubKey  = &sECDHPubKey{}
)

const (
	CCurveSize   = 256
	CECDHKeyType = "go-peer/ecdh"
)

/*
 * PRIVATE KEY
 */

type sECDHPrivKey struct {
	fPubKey  IEphPubKey
	fPrivKey *ecdh.PrivateKey
}

func newECDHPrivKey(pPrivKey *ecdh.PrivateKey) IEphPrivKey {
	return &sECDHPrivKey{
		fPubKey:  newECDHPubKey(pPrivKey.PublicKey()),
		fPrivKey: pPrivKey,
	}
}

func NewECDHPrivKey() IEphPrivKey {
	privKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return newECDHPrivKey(privKey)
}

func LoadECDHPrivKey(pPrivKey interface{}) IEphPrivKey {
	switch x := pPrivKey.(type) {
	case []byte:
		privKey, err := ecdh.P256().NewPrivateKey(x)
		if err != nil {
			return nil
		}
		return newECDHPrivKey(privKey)
	case string:
		x = skipSpaceChars(x)
		var (
			prefix = fmt.Sprintf(cPrivKeyPrefixTemplate, CECDHKeyType)
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
		return LoadECDHPrivKey(pbytes)
	default:
		panic("unsupported type")
	}
}

func (p *sECDHPrivKey) GetSharedKey(pEphPubKey IEphPubKey) ([]byte, error) {
	ecdhPubKey, ok := pEphPubKey.(*sECDHPubKey)
	if !ok {
		return nil, errors.NewError("invalid ecdh public key")
	}
	sharedKey, err := p.fPrivKey.ECDH(ecdhPubKey.fPubKey)
	if err != nil {
		return nil, errors.WrapError(err, "create shared key")
	}
	return sharedKey, nil
}

func (p *sECDHPrivKey) GetPubKey() IEphPubKey {
	return p.fPubKey
}

func (p *sECDHPrivKey) ToBytes() []byte {
	return p.fPrivKey.Bytes()
}

func (p *sECDHPrivKey) ToString() string {
	return fmt.Sprintf(cPrivKeyPrefixTemplate+"%X"+cKeySuffix, p.GetType(), p.ToBytes())
}

func (p *sECDHPrivKey) GetType() string {
	return CECDHKeyType
}

func (p *sECDHPrivKey) GetSize() uint64 {
	return CCurveSize
}

/*
 * PUBLIC KEY
 */

type sECDHPubKey struct {
	fPubKey *ecdh.PublicKey
}

func newECDHPubKey(pPubKey *ecdh.PublicKey) IEphPubKey {
	return &sECDHPubKey{
		fPubKey: pPubKey,
	}
}

func LoadECDHPubKey(pPubKey interface{}) IEphPubKey {
	switch x := pPubKey.(type) {
	case []byte:
		pub, err := ecdh.P256().NewPublicKey(x)
		if err != nil {
			return nil
		}
		return newECDHPubKey(pub)
	case string:
		x = skipSpaceChars(x)
		var (
			prefix = fmt.Sprintf(cPubKeyPrefixTemplate, CECDHKeyType)
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
		return LoadECDHPubKey(pbytes)
	default:
		panic("unsupported type")
	}
}

func (p *sECDHPubKey) ToBytes() []byte {
	return p.fPubKey.Bytes()
}

func (p *sECDHPubKey) ToString() string {
	return fmt.Sprintf(cPubKeyPrefixTemplate+"%X"+cKeySuffix, p.GetType(), p.ToBytes())
}

func (p *sECDHPubKey) GetType() string {
	return CECDHKeyType
}

func (p *sECDHPubKey) GetSize() uint64 {
	return CCurveSize
}
