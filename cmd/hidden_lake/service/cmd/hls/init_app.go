package main

import (
	"flag"
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/app"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/filesystem"
	"github.com/number571/go-peer/pkg/types"

	pkg_config "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

// initApp work with the raw data = read files, read args
func initApp() (types.ICommand, error) {
	var (
		inputPath string
		inputKey  string
	)

	flag.StringVar(&inputPath, "path", ".", "path to config/database files")
	flag.StringVar(&inputKey, "key", "", "input private key from file")
	flag.Parse()

	cfg, err := pkg_config.InitConfig(fmt.Sprintf("%s/%s", inputPath, pkg_settings.CPathCFG), nil)
	if err != nil {
		return nil, errors.WrapError(err, "init config")
	}

	var privKey asymmetric.IPrivKey
	switch inputKey {
	case "":
		privKey = asymmetric.NewRSAPrivKey(cfg.GetKeySizeBits())
	default:
		privKeyStr, err := filesystem.OpenFile(inputKey).Read()
		if err != nil {
			return nil, errors.WrapError(err, "read public key")
		}
		privKey = asymmetric.LoadRSAPrivKey(string(privKeyStr))
	}

	if privKey == nil {
		return nil, errors.NewError("private key is invalid")
	}

	return app.NewApp(cfg, privKey, inputPath), nil
}
