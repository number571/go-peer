package network

import (
	"math"

	"github.com/number571/go-peer/encoding"
)

var (
	_ iPackage = sPackage{}
)

type sPackage []byte

func newPackage(bytes []byte) iPackage {
	return sPackage(bytes)
}

// Size of package in big endian bytes.
func (pack sPackage) SizeToBytes() []byte {
	return encoding.Uint64ToBytes(uint64(pack.size()))
}

// From big endian bytes to uint size.
func (pack sPackage) BytesToSize() uint64 {
	// incorrect package
	if len(pack.bytes()) < cSizeUint {
		return math.MaxUint64
	}
	return encoding.BytesToUint64(pack.bytes())
}

// Size of package in big endian bytes.
func (pack sPackage) size() uint64 {
	return uint64(len(pack.bytes()))
}

// Bytes of package.
func (pack sPackage) bytes() []byte {
	return []byte(pack)
}
