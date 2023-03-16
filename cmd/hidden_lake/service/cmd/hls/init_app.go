package main

import (
	"flag"
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/app"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/filesystem"
	"github.com/number571/go-peer/pkg/types"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

// initApp work with the raw data = read files, read args
func initApp() (types.ICommand, error) {
	var (
		inputKey string
	)

	flag.StringVar(&inputKey, "key", "", "input private key from file")
	flag.Parse()

	var privKey asymmetric.IPrivKey
	switch inputKey {
	case "":
		privKey = asymmetric.NewRSAPrivKey(pkg_settings.CAKeySize)
	default:
		privKeyStr, err := filesystem.OpenFile(inputKey).Read()
		if err != nil {
			return nil, err
		}
		privKey = asymmetric.LoadRSAPrivKey(string(privKeyStr))
	}

	if privKey == nil {
		return nil, fmt.Errorf("private key is invalid")
	}

	cfg, err := initConfig()
	if err != nil {
		return nil, err
	}

	return app.NewApp(cfg, privKey), nil
}
