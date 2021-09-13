package crypto

import "bytes"

// Increase entropy by multiple hashing.
func RaiseEntropy(info, salt []byte, bits int) []byte {
	lim := uint64(1 << uint(bits))
	for i := uint64(0); i < lim; i++ {
		info = HashSum(bytes.Join(
			[][]byte{
				info,
				salt,
			},
			[]byte{},
		))
	}
	return info
}
