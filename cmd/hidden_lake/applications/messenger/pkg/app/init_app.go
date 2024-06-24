package app

import (
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	"github.com/number571/go-peer/internal/flag"
	"github.com/number571/go-peer/internal/initapp"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"
)

func InitApp(pArgs []string, pDefaultPath, pDefaultKey string) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", pDefaultPath), "/")

	cfg, err := config.InitConfig(filepath.Join(inputPath, settings.CPathYML), nil)
	if err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}

	// 1. https://crypto.stackexchange.com/questions/3288/is-truncating-a-hashed-private-key-with-sha-1-safe-to-use-as-the-symmetric-key-f
	// 2. https://www.reddit.com/r/crypto/comments/zwmoqf/can_a_private_key_be_used_for_symmetric_encryption/
	inputKey := flag.GetFlagValue(pArgs, "key", pDefaultKey)
	password, err := initapp.GetPassword(inputKey, settings.CServiceFullName)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetPassword, err)
	}

	return NewApp(cfg, password, inputPath), nil
}
