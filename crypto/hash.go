package crypto

import "crypto/sha256"

const (
	HashSize = sha256.Size
)

// Used SHA256.
func HashSum(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}
