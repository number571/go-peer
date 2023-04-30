package main

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/types"

	hls_app "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/app"
	hls_config "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

const (
	hlmURL = "localhost:9591"
)

func initApp(pPathTo string) (types.ICommand, types.ICommand, error) {
	privKey := asymmetric.NewRSAPrivKey(hls_settings.CAKeySize)
	if privKey == nil {
		return nil, nil, fmt.Errorf("private key is invalid")
	}

	cfgHLS, err := hls_config.InitConfig(
		fmt.Sprintf("%s/%s", pPathTo, hls_settings.CPathCFG),
		&hls_config.SConfig{
			FNetwork: "mobile_" + hls_settings.CServiceName,
			FAddress: &hls_config.SAddress{
				FTCP:  "localhost:9571",
				FHTTP: "localhost:9572",
			},
			FServices: map[string]string{
				settings.CTitlePattern: "localhost:9592",
			},
		},
	)
	if err != nil {
		return nil, nil, err
	}

	cfgHLM, err := config.InitConfig(
		fmt.Sprintf("%s/%s", pPathTo, settings.CPathCFG),
		&config.SConfig{
			FStorageKey: "mobile_" + settings.CServiceName,
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
		return nil, nil, err
	}

	return hls_app.NewApp(cfgHLS, privKey, pPathTo), app.NewApp(cfgHLM, pPathTo), nil
}
