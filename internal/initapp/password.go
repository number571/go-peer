package initapp

import (
	"os"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/utils"
)

func GetPassword(pPaswPath string) (string, error) {
	if _, err := os.Stat(pPaswPath); os.IsNotExist(err) {
		pasw := random.NewCSPRNG().GetString(32)
		if err := os.WriteFile(pPaswPath, []byte(pasw), 0o600); err != nil {
			return "", utils.MergeErrors(ErrWritePrivateKey, err)
		}
		return pasw, nil
	}
	paswBytes, err := os.ReadFile(pPaswPath)
	if err != nil {
		return "", utils.MergeErrors(ErrReadPrivateKey, err)
	}
	return string(paswBytes), nil
}
