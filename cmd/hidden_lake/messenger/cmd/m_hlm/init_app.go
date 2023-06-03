package main

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/types"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

const (
	hlmURL = "localhost:9591"
)

func initApp(pPathTo string) (types.ICommand, error) {
	privKey := asymmetric.NewRSAPrivKey(hls_settings.CAKeySize)
	if privKey == nil {
		return nil, errors.NewError("private key is invalid")
	}

	cfgHLM, err := config.InitConfig(
		fmt.Sprintf("%s/%s", pPathTo, settings.CPathCFG),
		&config.SConfig{
			FStorageKey: "mobile_" + pkg_settings.CServiceName,
			FAddress: &config.SAddress{
				FInterface: hlmURL,
				FIncoming:  "localhost:9592",
			},
			FConnection: &config.SConnection{
				FService: "localhost:9572",
			},
		},
	)
	if err != nil {
		return nil, errors.WrapError(err, "init config HLM")
	}

	return app.NewApp(cfgHLM, pPathTo), nil
}
