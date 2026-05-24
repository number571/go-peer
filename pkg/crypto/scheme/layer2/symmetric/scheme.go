package symmetric

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

const (
	// Salt + Nonce + AuthTag + Int32[len]
	// 16 + 12 + 16 + 4 = 48 additional bytes to origin message
	CMessageHeadSize = 0 +
		1*cSaltSize +
		1*cNonceSize +
		1*symmetric.CCipherBlockSize +
		1*encoding.CSizeUint32
)

const (
	cSaltSize  = 16 // salt for key generation
	cNonceSize = 12 // gcm nonce
)

func init() {
	scheme, err := NewScheme(128)
	if err != nil {
		panic(err)
	}
	if (scheme.GetMessageSize() - scheme.GetPayloadLimit()) != CMessageHeadSize {
		panic("incorrect calculated head size of message")
	}
}

var (
	_ layer2.IScheme = &sScheme{}
)

type sScheme struct {
	fMessageSize  uint64
	fPayloadLimit uint64
}

func NewScheme(pMessageSize uint64) (layer2.IScheme, error) {
	scheme := &sScheme{
		fMessageSize: pMessageSize,
	}

	encKey := make([]byte, symmetric.CCipherKeySize)
	cipher := symmetric.NewCipherGCM(encKey)
	encMsg := scheme.encryptWithPadding(cipher, []byte{}, 0)

	structSize := uint64(len(encMsg))
	if structSize >= pMessageSize {
		return nil, errors.Join(ErrStructGTEMessageSize, fmt.Errorf("struct size = %d", structSize)) //nolint:err113
	}

	scheme.fPayloadLimit = pMessageSize - structSize
	return scheme, nil
}

func (p *sScheme) GetRandomKey() layer2.IParticipantKey {
	return symmetric.NewCipherGCM(random.NewRandom().GetBytes(symmetric.CCipherKeySize))
}

func (p *sScheme) GetMessageSize() uint64 {
	return p.fMessageSize
}

func (p *sScheme) GetPayloadLimit() uint64 {
	return p.fPayloadLimit
}

func (p *sScheme) EncryptMessage(pKey layer2.IParticipantKey, pMsg []byte) ([]byte, error) {
	if _, ok := pKey.(symmetric.ICipher); !ok {
		return nil, ErrInvalidKeyType
	}

	var (
		payloadLimit = p.fPayloadLimit
		resultSize   = uint64(len(pMsg))
	)

	if resultSize > payloadLimit {
		return nil, ErrLimitMessageSize
	}

	return p.encryptWithPadding(pKey, pMsg, payloadLimit-resultSize), nil
}

func (p *sScheme) DecryptMessage(pKeysContainer layer2.IKeysContainer, pMsg []byte) (layer2.IParticipantKey, []byte, error) {
	if uint64(len(pMsg)) != p.fMessageSize {
		return nil, nil, ErrInvalidMessageSize
	}

	listKeys := pKeysContainer.List()
	salt := pMsg[:cSaltSize]

	var (
		decrypted bool
		resKey    layer2.IParticipantKey
		resMsg    []byte
	)

	for _, k := range listKeys {
		keyBuilder := keybuilder.NewKeyBuilder(0, salt)
		encKey := keyBuilder.Build(k.ToString(), symmetric.CCipherKeySize)

		cipher := symmetric.NewCipherGCM(encKey)
		decMsg := cipher.DecryptBytes(pMsg[cSaltSize:])
		if decMsg == nil {
			continue
		}

		lenb := [encoding.CSizeUint32]byte{}
		copy(lenb[:], decMsg[:encoding.CSizeUint32])
		msg := decMsg[encoding.CSizeUint32:]

		msgSize := encoding.BytesToUint32(lenb)
		if msgSize > uint32(len(msg)) { // nolint: gosec
			return nil, nil, ErrDecodeMessage
		}

		decrypted = true
		resKey = k
		resMsg = msg[:msgSize]
	}

	if decrypted {
		return resKey, resMsg, nil
	}
	return nil, nil, ErrDecryptMessage
}

func (p *sScheme) encryptWithPadding(
	pKey layer2.IParticipantKey,
	pMsg []byte,
	pPadd uint64,
) []byte {
	var (
		rand = random.NewRandom()
		salt = rand.GetBytes(cSaltSize)
	)

	lenb := encoding.Uint32ToBytes(uint32(len(pMsg))) // nolint: gosec
	data := bytes.Join([][]byte{lenb[:], pMsg, rand.GetBytes(pPadd)}, []byte{})

	keyBuilder := keybuilder.NewKeyBuilder(0, salt)
	encKey := keyBuilder.Build(pKey.ToString(), symmetric.CCipherKeySize)
	return bytes.Join([][]byte{
		salt,
		symmetric.NewCipherGCM(encKey).EncryptBytes(data),
	}, []byte{})
}
