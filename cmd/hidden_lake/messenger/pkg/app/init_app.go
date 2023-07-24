package app

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	"github.com/number571/go-peer/internal/flag"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/types"
)

func InitApp() (types.ICommand, error) {
	inputPath := flag.GetFlagValue("path", ".")

	cfg, err := config.InitConfig(fmt.Sprintf("%s/%s", inputPath, settings.CPathCFG), nil)
	if err != nil {
		return nil, errors.WrapError(err, "init config")
	}

	return NewApp(cfg, inputPath), nil
}
