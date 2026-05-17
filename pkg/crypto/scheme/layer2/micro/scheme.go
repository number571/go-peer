package micro

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	hybrid "github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
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

	tmpCipher := symmetric.NewCipher(random.NewRandom().GetBytes(symmetric.CCipherKeySize))
	encMsg, err := scheme.encryptWithPadding(tmpCipher, []byte{}, 0)
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
	return symmetric.NewCipher(random.NewRandom().GetBytes(symmetric.CCipherKeySize))
}

func (p *sScheme) GetMessageSize() uint64 {
	return p.fMessageSize
}

func (p *sScheme) GetPayloadLimit() uint64 {
	return p.fPayloadLimit
}

func (p *sScheme) EncryptMessage(pRecv hybrid.IParticipantKey, pMsg []byte) ([]byte, error) {
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
	return p.encryptWithPadding(recv, pMsg, payloadLimit-resultSize)
}

func (p *sScheme) DecryptMessage(pListKeys hybrid.IKeysContainer, pMsg []byte) (hybrid.IParticipantKey, []byte, error) {
	if uint64(len(pMsg)) != p.fMessageSize {
		return nil, nil, ErrMessageSize
	}
	listCiphers, ok := pListKeys.(symmetric.IListCiphers)
	if !ok {
		return nil, nil, ErrInvalidKeyType
	}
	list := listCiphers.Get()
	for _, c := range list {
		dec := c.DecryptBytes(pMsg)
		if len(dec) <= cSaltSize+hashing.CHasherSize+encoding.CSizeUint32 {
			continue
		}
		var (
			salt = dec[:cSaltSize]
			hmac = dec[cSaltSize : cSaltSize+hashing.CHasherSize]
			data = dec[cSaltSize+hashing.CHasherSize:]
		)
		check := hashing.NewHMACHasher(c.ToBytes(), bytes.Join(
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
		return c, msg[:msgSize], nil
	}
	return nil, nil, ErrDecryptMessage
}

func (p *sScheme) encryptWithPadding(
	pCipher symmetric.ICipher,
	pMsg []byte,
	pPadd uint64,
) ([]byte, error) {
	var (
		rand = random.NewRandom()
		salt = rand.GetBytes(cSaltSize)
	)
	lenb := encoding.Uint32ToBytes(uint32(len(pMsg))) // nolint: gosec
	data := bytes.Join([][]byte{lenb[:], pMsg, rand.GetBytes(pPadd)}, []byte{})
	return pCipher.EncryptBytes(bytes.Join(
		[][]byte{
			salt,
			hashing.NewHMACHasher(pCipher.ToBytes(), bytes.Join(
				[][]byte{salt, data},
				[]byte{},
			)).ToBytes(),
			data,
		},
		[]byte{},
	)), nil
}
