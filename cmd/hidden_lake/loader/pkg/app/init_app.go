package app

import (
	"fmt"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/loader/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/loader/pkg/settings"
	"github.com/number571/go-peer/internal/flag"
	"github.com/number571/go-peer/pkg/types"
)

func InitApp(pDefaultPath string) (types.IApp, error) {
	inputPath := strings.TrimSuffix(flag.GetFlagValue("path", pDefaultPath), "/")

	cfg, err := config.InitConfig(fmt.Sprintf("%s/%s", inputPath, settings.CPathYML), nil)
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	return NewApp(cfg, inputPath), nil
}
