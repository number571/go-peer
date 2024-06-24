package initapp

import (
	"os"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/utils"
)

func GetPassword(pKeyPath, pServiceName string) (string, error) {
	keyBytes, err := os.ReadFile(pKeyPath)
	if err != nil {
		return "", utils.MergeErrors(ErrReadPrivateKey, err)
	}
	return hashing.NewHMACSHA256Hasher(
		keyBytes,
		[]byte(pServiceName),
	).ToString(), nil
}
