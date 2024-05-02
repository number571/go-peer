package payload

import (
	"bytes"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IPayload64 = sPayload64{}
)

type sPayload64 []byte

func NewPayload64(pHead uint64, pData []byte) IPayload64 {
	bHead := encoding.Uint64ToBytes(pHead)
	return sPayload64(bytes.Join([][]byte{
		bHead[:],
		pData,
	}, []byte{}))
}

func LoadPayload64(pPayloadBytes []byte) IPayload64 {
	if len(pPayloadBytes) < encoding.CSizeUint64 {
		return nil
	}
	return sPayload64(pPayloadBytes)
}

func (p sPayload64) GetHead() uint64 {
	bHead := [encoding.CSizeUint64]byte{}
	copy(bHead[:], p[:encoding.CSizeUint64])
	return encoding.BytesToUint64(bHead)
}

func (p sPayload64) GetBody() []byte {
	return p[encoding.CSizeUint64:]
}

func (p sPayload64) ToBytes() []byte {
	return p[:]
}
