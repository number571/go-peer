package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
	"github.com/number571/go-peer/internal/flag"
	"github.com/number571/go-peer/pkg/types"
)

func InitApp(pArgs []string, pDefaultPath string) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", pDefaultPath), "/")

	cfg, err := config.InitConfig(filepath.Join(inputPath, settings.CPathYML), nil)
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	if err := os.MkdirAll(settings.CPathSTG, 0o777); err != nil {
		return nil, fmt.Errorf("init storage: %w", err)
	}

	return NewApp(cfg, inputPath), nil
}
