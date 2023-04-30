package main

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/app"
	"github.com/number571/go-peer/pkg/types"

	pkg_config "github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

// initApp work with the raw data = read files, read args
func initApp(path string) (types.ICommand, error) {
	cfg, err := pkg_config.InitConfig(
		fmt.Sprintf("%s/%s", path, pkg_settings.CPathCFG),
		&pkg_config.SConfig{
			FAddress:    "localhost:9581",
			FConnection: "localhost:9571",
		},
	)
	if err != nil {
		return nil, err
	}
	return app.NewApp(cfg, path), nil
}
