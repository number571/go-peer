package payload

import (
	"bytes"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IPayload32 = sPayload32{}
)

type sPayload32 []byte

func NewPayload32(pHead uint32, pData []byte) IPayload32 {
	bHead := encoding.Uint32ToBytes(pHead)
	return sPayload32(bytes.Join([][]byte{
		bHead[:],
		pData,
	}, []byte{}))
}

func LoadPayload32(pPayloadBytes []byte) IPayload32 {
	if len(pPayloadBytes) < encoding.CSizeUint32 {
		return nil
	}
	return sPayload32(pPayloadBytes)
}

func (p sPayload32) GetHead() uint32 {
	bHead := [encoding.CSizeUint32]byte{}
	copy(bHead[:], p[:encoding.CSizeUint32])
	return encoding.BytesToUint32(bHead)
}

func (p sPayload32) GetBody() []byte {
	return p[encoding.CSizeUint32:]
}

func (p sPayload32) ToBytes() []byte {
	return p[:]
}
