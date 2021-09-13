package crypto

import "crypto/sha256"

// Used SHA256.
func HashSum(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}
