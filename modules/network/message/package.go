package message

import (
	"bytes"

	"github.com/number571/go-peer/modules/encoding"
)

var (
	_ IPackage = &sPackage{}
)

type sPackage struct {
	fSize uint64
	fData []byte
}

func NewPackage(data []byte) IPackage {
	return &sPackage{
		fSize: uint64(len(data)),
		fData: data,
	}
}

func LoadPackage(packBytes []byte) IPackage {
	lenBytes := uint64(len(packBytes))
	if lenBytes < encoding.CSizeUint64 {
		return nil
	}

	res := [encoding.CSizeUint64]byte{}
	copy(res[:], packBytes[:encoding.CSizeUint64])

	packSize := encoding.BytesToUint64(res)
	if packSize != lenBytes-encoding.CSizeUint64 {
		return nil
	}

	return &sPackage{
		fSize: packSize,
		fData: packBytes[encoding.CSizeUint64:],
	}
}

func (pack *sPackage) Size() uint64 {
	return pack.fSize
}

func (pack *sPackage) Data() []byte {
	return pack.fData
}

func (pack *sPackage) ToBytes() []byte {
	res := encoding.Uint64ToBytes(pack.fSize)
	return bytes.Join([][]byte{
		res[:],
		pack.fData,
	}, []byte{})
}
