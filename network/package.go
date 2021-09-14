package network

import (
	"encoding/json"

	"github.com/number571/gopeer/encoding"
)

type Package []byte

// Size of package in big endian bytes.
func (pack Package) Size() []byte {
	length := uint64(len(pack.Bytes()))
	return encoding.Uint64ToBytes(length)
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
