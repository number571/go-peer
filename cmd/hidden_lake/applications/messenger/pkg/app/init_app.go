package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	"github.com/number571/go-peer/internal/flag"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"
)

func InitApp(pArgs []string, pDefaultPath, pDefaultPasw string) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", pDefaultPath), "/")

	cfg, err := config.InitConfig(filepath.Join(inputPath, settings.CPathYML), nil)
	if err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}

	inputPasw := flag.GetFlagValue(pArgs, "pasw", pDefaultPasw)
	password, err := getPassword(inputPasw)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetPassword, err)
	}

	return NewApp(cfg, password, inputPath), nil
}

func getPassword(pPaswPath string) (string, error) {
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
