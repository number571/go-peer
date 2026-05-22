package symmetric

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ layer2.IScheme = &sScheme{}
)

const (
	cSaltSize = 16 // salt for key generation
)

type sScheme struct {
	fMessageSize  uint64
	fPayloadLimit uint64
}

func NewScheme(pMessageSize uint64) layer2.IScheme {
	scheme := &sScheme{
		fMessageSize: pMessageSize,
	}

	encKey := make([]byte, symmetric.CCipherKeySize)
	cipher := symmetric.NewCipherGCM(encKey)
	encMsg := scheme.encryptWithPadding(cipher, []byte{}, 0)

	structSize := uint64(len(encMsg))
	if pMessageSize <= structSize {
		panic("the payload size is lower than struct size")
	}

	scheme.fPayloadLimit = pMessageSize - structSize
	return scheme
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

func (p *sScheme) EncryptMessage(pRecv layer2.IParticipantKey, pMsg []byte) ([]byte, error) {
	recv, ok := pRecv.(symmetric.ICipher)
	if !ok {
		return nil, ErrInvalidKeyType
	}
	var (
		payloadLimit = p.fPayloadLimit
		resultSize   = uint64(len(pMsg))
	)
	if resultSize > payloadLimit {
		return nil, ErrLimitMessageSize
	}
	return p.encryptWithPadding(recv, pMsg, payloadLimit-resultSize), nil
}

func (p *sScheme) DecryptMessage(pListKeys layer2.IKeysContainer, pMsg []byte) (layer2.IParticipantKey, []byte, error) {
	listCiphers, ok := pListKeys.(symmetric.IListCiphers)
	if !ok {
		return nil, nil, ErrInvalidKeyType
	}
	if uint64(len(pMsg)) != p.fMessageSize {
		return nil, nil, ErrInvalidMessageSize
	}

	salt := pMsg[:cSaltSize]
	list := listCiphers.Get()

	for _, c := range list {
		keyBuilder := keybuilder.NewKeyBuilder(0, salt)
		encKey := keyBuilder.Build(c.ToString(), symmetric.CCipherKeySize)

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

		return c, msg[:msgSize], nil
	}

	return nil, nil, ErrDecryptMessage
}

func (p *sScheme) encryptWithPadding(
	pCipher symmetric.ICipher,
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
	encKey := keyBuilder.Build(pCipher.ToString(), symmetric.CCipherKeySize)
	return bytes.Join([][]byte{
		salt,
		symmetric.NewCipherGCM(encKey).EncryptBytes(data),
	}, []byte{})
}
