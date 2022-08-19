package payload

import (
	"bytes"

	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/settings"
)

var (
	_ IPayload = sPayload{}
)

type sPayload []byte

func NewPayload(head uint64, data []byte) IPayload {
	return sPayload(bytes.Join([][]byte{
		encoding.Uint64ToBytes(head),
		data,
	}, []byte{}))
}

func LoadPayload(payloadBytes []byte) IPayload {
	if len(payloadBytes) < settings.CSizeUint64 {
		return nil
	}
	return sPayload(payloadBytes)
}

func (payload sPayload) Head() uint64 {
	return encoding.BytesToUint64(payload[:settings.CSizeUint64])
}

func (payload sPayload) Body() []byte {
	return payload[settings.CSizeUint64:]
}

func (payload sPayload) Bytes() []byte {
	return payload[:]
}
