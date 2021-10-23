package local

import (
	"encoding/json"

	"github.com/number571/gopeer/encoding"
)

type Package []byte

// Size of package in big endian bytes.
func (pack Package) Size() uint {
	return uint(len(pack.Bytes()))
}

// Size of package in big endian bytes.
func (pack Package) SizeToBytes() []byte {
	return encoding.Uint64ToBytes(uint64(pack.Size()))
}

// From big endian bytes to uint size.
func (pack Package) BytesToSize() uint {
	return uint(encoding.BytesToUint64(pack.Bytes()))
}

// Bytes of package.
func (pack Package) Bytes() []byte {
	return []byte(pack)
}

// Deserialize with JSON format.
func (pack Package) Deserialize() *Message {
	var msg = new(Message)
	err := json.Unmarshal(pack, msg)
	if err != nil {
		return nil
	}
	return msg
}
