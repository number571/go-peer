package main

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/app"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/types"

	pkg_config "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

const (
	hlsURL = "localhost:9572"
)

// initApp work with the raw data = read files, read args
func initApp(pPathTo string) (types.ICommand, error) {
	privKey := asymmetric.NewRSAPrivKey(pkg_settings.CAKeySize)
	if privKey == nil {
		return nil, errors.NewError("private key is invalid")
	}

	cfg, err := pkg_config.InitConfig(
		fmt.Sprintf("%s/%s", pPathTo, pkg_settings.CPathCFG),
		&pkg_config.SConfig{
			FNetwork: "mobile_" + pkg_settings.CServiceName,
			FAddress: &pkg_config.SAddress{
				FTCP:  "localhost:9571",
				FHTTP: hlsURL,
			},
		},
	)
	if err != nil {
		return nil, errors.WrapError(err, "init config")
	}

	return app.NewApp(cfg, privKey, pPathTo), nil
}
