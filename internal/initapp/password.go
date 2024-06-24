package initapp

import (
	"os"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/utils"
)

/*
1. https://crypto.stackexchange.com/questions/3288/is-truncating-a-hashed-private-key-with-sha-1-safe-to-use-as-the-symmetric-key-f
2. https://www.reddit.com/r/crypto/comments/zwmoqf/can_a_private_key_be_used_for_symmetric_encryption/
*/
func GetPrivKeyAsPassword(pKeyPath, pServiceName string) (string, error) {
	privKeyStr, err := os.ReadFile(pKeyPath)
	if err != nil {
		return "", utils.MergeErrors(ErrReadPrivateKey, err)
	}
	privKey := asymmetric.LoadRSAPrivKey(string(privKeyStr))
	if privKey == nil {
		return "", ErrInvalidPrivateKey
	}
	return hashing.NewHMACSHA256Hasher(
		privKey.ToBytes(),
		[]byte(pServiceName),
	).ToString(), nil
}
