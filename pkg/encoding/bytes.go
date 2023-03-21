package encoding

import (
	"bytes"
	"encoding/binary"
)

const (
	CSizeUint64 = 8 // bytes
	// Visual code's extension "Go v0.35.2" is piece of shit and can't
	// use CSizeUint64 constant for return [CSizeUint64]byte type
	// but can return with static constant cSizeUint64
	cSizeUint64 = CSizeUint64
)

// Uint64 to slice of bytes by big endian.
func Uint64ToBytes(pNum uint64) [cSizeUint64]byte {
	res := [CSizeUint64]byte{}

	var data = new(bytes.Buffer)
	err := binary.Write(data, binary.BigEndian, pNum)
	if err != nil {
		panic(err)
	}

	copy(res[:], data.Bytes()[:])
	return res
}

// Slice of bytes to uint64 by big endian.
func BytesToUint64(pBytes [cSizeUint64]byte) uint64 {
	return binary.BigEndian.Uint64(pBytes[:])
}
