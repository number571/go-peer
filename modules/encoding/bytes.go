package encoding

import (
	"bytes"
	"encoding/binary"
)

const (
	CSizeUint64 = 8 // bytes
)

// Uint64 to slice of bytes by big endian.
func Uint64ToBytes(num uint64) [CSizeUint64]byte {
	res := [CSizeUint64]byte{}

	var data = new(bytes.Buffer)
	err := binary.Write(data, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}

	copy(res[:], data.Bytes()[:])
	return res
}

// Slice of bytes to uint64 by big endian.
func BytesToUint64(bytes [CSizeUint64]byte) uint64 {
	return binary.BigEndian.Uint64(bytes[:])
}
