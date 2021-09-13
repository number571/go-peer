package encoding

import (
	"bytes"
	"encoding/binary"
)

// Uint64 to slice of bytes by big endian.
func ToBytes(num uint64) []byte {
	var data = new(bytes.Buffer)
	err := binary.Write(data, binary.BigEndian, num)
	if err != nil {
		return nil
	}
	return data.Bytes()
}
