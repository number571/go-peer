package encoding

import (
	"bytes"
	"encoding/binary"

	"github.com/number571/go-peer/settings"
)

// Uint64 to slice of bytes by big endian.
func Uint64ToBytes(num uint64) [settings.CSizeUint64]byte {
	res := [settings.CSizeUint64]byte{}

	var data = new(bytes.Buffer)
	err := binary.Write(data, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}

	copy(res[:], data.Bytes()[:])
	return res
}

// Slice of bytes to uint64 by big endian.
func BytesToUint64(bytes [settings.CSizeUint64]byte) uint64 {
	return binary.BigEndian.Uint64(bytes[:])
}
