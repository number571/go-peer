package crypto

import "crypto/rand"

// Generates a cryptographically strong pseudo-random sequence.
func Rand(max uint) []byte {
	var slice []byte = make([]byte, max)
	_, err := rand.Read(slice)
	if err != nil {
		return nil
	}
	return slice
}
