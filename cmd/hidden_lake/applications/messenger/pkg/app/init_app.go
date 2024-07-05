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

func InitApp(pArgs []string, pDefaultPath, pDefaultPaswPath string) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", pDefaultPath), "/")

	cfg, err := config.InitConfig(filepath.Join(inputPath, settings.CPathYML), nil)
	if err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}

	paswPath := flag.GetFlagValue(pArgs, "pasw", pDefaultPaswPath)
	password, err := initapp.GetPassword(paswPath)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetPassword, err)
	}

	return NewApp(cfg, password, inputPath), nil
}
