package initapp

import (
	"os"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/utils"
)

func GetPassword(pPaswPath string) (string, error) {
	if _, err := os.Stat(pPaswPath); os.IsNotExist(err) {
		password := random.NewCSPRNG().GetString(32)
		if err := os.WriteFile(pPaswPath, []byte(password), 0o600); err != nil {
			return "", utils.MergeErrors(ErrWritePassword, err)
		}
		return password, nil
	}

	password, err := os.ReadFile(pPaswPath)
	if err != nil {
		return "", utils.MergeErrors(ErrReadPassword, err)
	}
	return string(password), nil
}
