package app

import (
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	"github.com/number571/go-peer/internal/flag"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"
)

func InitApp(pArgs []string, pDefaultPath string) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", pDefaultPath), "/")

	cfgPath := filepath.Join(inputPath, settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil)
	if err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}

	return NewApp(cfg, inputPath), nil
}
