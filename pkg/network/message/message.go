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
	fHash    []byte
	fPayload payload.IPayload
}

func NewMessage(pPld payload.IPayload, pKey []byte) IMessage {
	return &sMessage{
		fHash: hashing.NewHMACSHA256Hasher(
			pKey,
			pPld.ToBytes(),
		).ToBytes(),
		fPayload: pPld,
	}
}

func LoadMessage(pData, pKey []byte) IMessage {
	// check Hash[uN]
	if len(pData) < hashing.CSHA256Size {
		return nil
	}

	hashRecv := pData[:hashing.CSHA256Size]
	payloadBytes := pData[hashing.CSHA256Size:]
	if !bytes.Equal(
		hashRecv,
		hashing.NewHMACSHA256Hasher(
			pKey,
			payloadBytes,
		).ToBytes(),
	) {
		return nil
	}

	// check Head[u64]
	pld := payload.LoadPayload(payloadBytes)
	if pld == nil {
		return nil
	}

	return &sMessage{
		fHash:    hashRecv,
		fPayload: pld,
	}
}

func (p *sMessage) GetHash() []byte {
	return p.fHash
}

func (p *sMessage) GetPayload() payload.IPayload {
	return p.fPayload
}

func (p *sMessage) GetBytes() []byte {
	return bytes.Join(
		[][]byte{
			p.fHash,
			p.fPayload.ToBytes(),
		},
		[]byte{},
	)
}
