package main

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/app"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/types"

	pkg_config "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

// initApp work with the raw data = read files, read args
func initApp(path string) (types.ICommand, error) {
	privKey := asymmetric.NewRSAPrivKey(pkg_settings.CAKeySize)
	if privKey == nil {
		return nil, fmt.Errorf("private key is invalid")
	}

	cfg, err := pkg_config.InitConfig(
		fmt.Sprintf("%s/%s", path, pkg_settings.CPathCFG),
		&pkg_config.SConfig{
			FNetwork: "android_" + pkg_settings.CServiceName,
			FAddress: &pkg_config.SAddress{
				FTCP:  "localhost:9571",
				FHTTP: "localhost:9572",
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return app.NewApp(cfg, privKey, path), nil
}
