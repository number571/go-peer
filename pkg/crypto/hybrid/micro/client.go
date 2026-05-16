package micro

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/hybrid"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/types"
)

var (
	_ hybrid.IScheme = &sScheme{}
)

const (
	cSaltSize = 16 // bytes
)

type sScheme struct {
	fMessageSize  uint64
	fPayloadLimit uint64
}

func NewScheme(pMessageSize uint64) hybrid.IScheme {
	scheme := &sScheme{
		fMessageSize: pMessageSize,
	}

	tmpKey := random.NewRandom().GetBytes(symmetric.CCipherKeySize)
	encMsg, err := scheme.encryptWithPadding(tmpKey, []byte{}, 0)
	if err != nil {
		panic(err)
	}

	structSize := uint64(len(encMsg))
	if pMessageSize <= structSize {
		panic("the payload size is lower than struct size")
	}

	scheme.fPayloadLimit = pMessageSize - structSize
	return scheme
}

func (p *sScheme) GetRandomKey() hybrid.IParticipantKey {
	return types.NewConverter(random.NewRandom().GetBytes(symmetric.CCipherKeySize))
}

func (p *sScheme) GetMessageSize() uint64 {
	return p.fMessageSize
}

func (p *sScheme) GetPayloadLimit() uint64 {
	return p.fPayloadLimit
}

func (p *sScheme) EncryptMessage(pKey hybrid.IParticipantKey, pMsg []byte) ([]byte, error) {
	recv := pKey.ToBytes()
	if len(recv) != symmetric.CCipherKeySize {
		return nil, ErrInvalidKeyType
	}
	var (
		payloadLimit = p.fPayloadLimit
		resultSize   = uint64(len(pMsg))
	)
	if resultSize > payloadLimit {
		return nil, ErrLimitMessageSize
	}
	return p.encryptWithPadding(recv, pMsg, payloadLimit-resultSize)
}

func (p *sScheme) DecryptMessage(pListKeys hybrid.IKeysContainer, pMsg []byte) (hybrid.IParticipantKey, []byte, error) {
	listKeys, ok := pListKeys.([]hybrid.IParticipantKey)
	if !ok {
		return nil, nil, ErrInvalidKeyType
	}
	if uint64(len(pMsg)) != p.fMessageSize {
		return nil, nil, ErrMessageSize
	}
	for _, k := range listKeys {
		cipher := symmetric.NewCipher(k.ToBytes())
		if cipher == nil {
			continue
		}
		dec := cipher.DecryptBytes(pMsg)
		if len(dec) <= cSaltSize+hashing.CHasherSize+encoding.CSizeUint32 {
			continue
		}
		var (
			salt = dec[:cSaltSize]
			hmac = dec[cSaltSize : cSaltSize+hashing.CHasherSize]
			data = dec[cSaltSize+hashing.CHasherSize:]
		)
		check := hashing.NewHMACHasher(k.ToBytes(), bytes.Join(
			[][]byte{salt, data},
			[]byte{},
		)).ToBytes()
		if !bytes.Equal(check, hmac) {
			continue
		}
		lenb := [encoding.CSizeUint32]byte{}
		copy(lenb[:], data[:encoding.CSizeUint32])
		msg := data[encoding.CSizeUint32:]
		msgSize := encoding.BytesToUint32(lenb)
		if msgSize > uint32(len(msg)) { // nolint: gosec
			return nil, nil, ErrDecodeMessage
		}
		return k, msg[:msgSize], nil
	}
	return nil, nil, ErrDecryptMessage
}

func (p *sScheme) encryptWithPadding(
	pKey []byte,
	pMsg []byte,
	pPadd uint64,
) ([]byte, error) {
	var (
		rand = random.NewRandom()
		salt = rand.GetBytes(cSaltSize)
	)
	lenb := encoding.Uint32ToBytes(uint32(len(pMsg))) // nolint: gosec
	data := bytes.Join([][]byte{lenb[:], pMsg, rand.GetBytes(pPadd)}, []byte{})
	cipher := symmetric.NewCipher(pKey)
	return cipher.EncryptBytes(bytes.Join(
		[][]byte{
			salt,
			hashing.NewHMACHasher(pKey, bytes.Join(
				[][]byte{salt, data},
				[]byte{},
			)).ToBytes(),
			data,
		},
		[]byte{},
	)), nil
}
