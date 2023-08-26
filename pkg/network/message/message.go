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

func NewMessage(pPld payload.IPayload) IMessage {
	return &sMessage{
		fPayload:     pPld,
		fHashPayload: getHash(pPld.ToBytes()),
	}
}

func LoadMessage(pData []byte) IMessage {
	// check Hash[uN]
	if len(pData) < hashing.CSHA256Size {
		return nil
	}

	hashRecv := pData[:hashing.CSHA256Size]
	pldBytes := pData[hashing.CSHA256Size:]

	if !bytes.Equal(
		hashRecv,
		getHash(pldBytes),
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

func (p *sMessage) ToBytes() []byte {
	return bytes.Join(
		[][]byte{
			p.fHashPayload,
			p.fPayload.ToBytes(),
		},
		[]byte{},
	)
}

func getHash(pBytes []byte) []byte {
	return hashing.NewSHA256Hasher(pBytes).ToBytes()
}
