package encoding

import (
	"bytes"
	"encoding/binary"

	"github.com/number571/go-peer/settings"
)

// Uint64 to slice of bytes by big endian.
func Uint64ToBytes(num uint64) []byte {
	var data = new(bytes.Buffer)
	err := binary.Write(data, binary.BigEndian, num)
	if err != nil {
		return nil
	}
	return data.Bytes()
}

// Slice of bytes to uint64 by big endian.
func BytesToUint64(bytes []byte) uint64 {
	if len(bytes) != settings.CSizeUint64 {
		panic("len(bytes) != settings.CSizeUint64")
	}
	return binary.BigEndian.Uint64(bytes)
}
