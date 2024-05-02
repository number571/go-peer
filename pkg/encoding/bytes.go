package encoding

import (
	"bytes"
	"encoding/binary"
)

const (
	CSizeUint32 = 4 // bytes
	CSizeUint64 = 8 // bytes
)

// Uint64 to slice of bytes by big endian.
func Uint64ToBytes(pNum uint64) [CSizeUint64]byte {
	res := [CSizeUint64]byte{}

	var data = new(bytes.Buffer)
	_ = binary.Write(data, binary.BigEndian, pNum)

	copy(res[:], data.Bytes())
	return res
}

// Uint32 to slice of bytes by big endian.
func Uint32ToBytes(pNum uint32) [CSizeUint32]byte {
	res := [CSizeUint32]byte{}

	var data = new(bytes.Buffer)
	_ = binary.Write(data, binary.BigEndian, pNum)

	copy(res[:], data.Bytes())
	return res
}

// Slice of bytes to uint64 by big endian.
func BytesToUint64(pBytes [CSizeUint64]byte) uint64 {
	return binary.BigEndian.Uint64(pBytes[:])
}

// Slice of bytes to uint32 by big endian.
func BytesToUint32(pBytes [CSizeUint32]byte) uint32 {
	return binary.BigEndian.Uint32(pBytes[:])
}
