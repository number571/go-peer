package local

import (
	"encoding/json"

	"github.com/number571/go-peer/encoding"
)

var (
	_ IPackage = sPackage{}
)

type sPackage []byte

func LoadPackage(data []byte) IPackage {
	return sPackage(data)
}

// Size of package in big endian bytes.
func (pack sPackage) Size() uint64 {
	return uint64(len(pack.Bytes()))
}

// Bytes of package.
func (pack sPackage) Bytes() []byte {
	return []byte(pack)
}

// Size of package in big endian bytes.
func (pack sPackage) SizeToBytes() []byte {
	return encoding.Uint64ToBytes(uint64(pack.Size()))
}

// From big endian bytes to uint size.
func (pack sPackage) BytesToSize() uint64 {
	return encoding.BytesToUint64(pack.Bytes())
}

// Deserialize with JSON format.
func (pack sPackage) ToMessage() IMessage {
	var msg = new(sMessage)
	err := json.Unmarshal(pack, msg)
	if err != nil {
		return nil
	}
	return msg
}
