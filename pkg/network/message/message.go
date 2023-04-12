package message

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/payload"
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	fPayload     payload.IPayload
	fHashPayload []byte
}

func NewMessage(pPld payload.IPayload, pKey []byte) IMessage {
	return &sMessage{
		fPayload:     pPld,
		fHashPayload: getHash(pKey, pPld.ToBytes()),
	}
}

func LoadMessage(pData, pKey []byte) IMessage {
	// check Hash[uN]
	if len(pData) < hashing.CSHA256Size {
		return nil
	}

	hashRecv := pData[:hashing.CSHA256Size]
	pldBytes := pData[hashing.CSHA256Size:]

	if !bytes.Equal(
		hashRecv,
		getHash(pKey, pldBytes),
	) {
		return nil
	}

	// check Head[u64]
	pld := payload.LoadPayload(pldBytes)
	if pld == nil {
		return nil
	}

	return &sMessage{
		fPayload:     pld,
		fHashPayload: hashRecv,
	}
}

func (p *sMessage) GetHash() []byte {
	return p.fHashPayload
}

func (p *sMessage) GetPayload() payload.IPayload {
	return p.fPayload
}

func (p *sMessage) GetBytes() []byte {
	return bytes.Join(
		[][]byte{
			p.fHashPayload,
			p.fPayload.ToBytes(),
		},
		[]byte{},
	)
}

func getHash(pkey, pBytes []byte) []byte {
	rawKey := bytes.Join(
		[][]byte{[]byte("__"), pkey, []byte("__")},
		[]byte{},
	)
	return hashing.NewHMACSHA256Hasher(
		hashing.NewSHA256Hasher(rawKey).ToBytes(),
		pBytes,
	).ToBytes()
}
