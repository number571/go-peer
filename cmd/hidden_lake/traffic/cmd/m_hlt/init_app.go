package main

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/app"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/types"

	pkg_config "github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

const (
	hltURL = "localhost:9581"
)

// initApp work with the raw data = read files, read args
func initApp(pPathTo string) (types.ICommand, error) {
	cfg, err := pkg_config.InitConfig(
		fmt.Sprintf("%s/%s", pPathTo, pkg_settings.CPathCFG),
		&pkg_config.SConfig{
			FAddress:    hltURL,
			FConnection: "localhost:9571",
		},
	)
	if err != nil {
		return nil, errors.WrapError(err, "init config")
	}
	return app.NewApp(cfg, pPathTo), nil
}
