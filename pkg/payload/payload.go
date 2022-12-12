package payload

import (
	"bytes"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IPayload = sPayload{}
)

type sPayload []byte

func NewPayload(head uint64, data []byte) IPayload {
	bHead := encoding.Uint64ToBytes(head)
	return sPayload(bytes.Join([][]byte{
		bHead[:],
		data,
	}, []byte{}))
}

func LoadPayload(payloadBytes []byte) IPayload {
	if len(payloadBytes) < encoding.CSizeUint64 {
		return nil
	}
	return sPayload(payloadBytes)
}

func (payload sPayload) Head() uint64 {
	bHead := [encoding.CSizeUint64]byte{}
	copy(bHead[:], payload[:encoding.CSizeUint64])
	return encoding.BytesToUint64(bHead)
}

func (payload sPayload) Body() []byte {
	return payload[encoding.CSizeUint64:]
}

func (payload sPayload) Bytes() []byte {
	return payload[:]
}
