package payload

import (
	"bytes"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IPayload = sPayload{}
)

type sPayload []byte

func NewPayload(pHead uint64, pData []byte) IPayload {
	bHead := encoding.Uint64ToBytes(pHead)
	return sPayload(bytes.Join([][]byte{
		bHead[:],
		pData,
	}, []byte{}))
}

func LoadPayload(pPayloadBytes []byte) IPayload {
	if len(pPayloadBytes) < encoding.CSizeUint64 {
		return nil
	}
	return sPayload(pPayloadBytes)
}

func (p sPayload) GetHead() uint64 {
	bHead := [encoding.CSizeUint64]byte{}
	copy(bHead[:], p[:encoding.CSizeUint64])
	return encoding.BytesToUint64(bHead)
}

func (p sPayload) GetBody() []byte {
	return p[encoding.CSizeUint64:]
}

func (p sPayload) ToBytes() []byte {
	return p[:]
}
