package app

import (
	"fmt"
	"strings"

	"github.com/number571/go-peer/internal/flag"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/filesystem"
	"github.com/number571/go-peer/pkg/types"

	pkg_config "github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

// initApp work with the raw data = read files, read args
func InitApp(pDefaultPath, pDefaultKey string) (types.ICommand, error) {
	inputPath := strings.TrimSuffix(flag.GetFlagValue("path", pDefaultPath), "/")
	inputKey := flag.GetFlagValue("key", pDefaultKey)

	cfg, err := pkg_config.InitConfig(fmt.Sprintf("%s/%s", inputPath, pkg_settings.CPathCFG), nil)
	if err != nil {
		return nil, errors.WrapError(err, "init config")
	}

	privKey, err := getPrivKey(cfg, inputKey)
	if err != nil {
		return nil, errors.WrapError(err, "get private key")
	}

	if privKey.GetSize() != cfg.GetSettings().GetKeySizeBits() {
		return nil, errors.NewError("size of private key is invalid")
	}

	return NewApp(cfg, privKey, inputPath), nil
}

func getPrivKey(pCfg pkg_config.IConfig, pKeyPath string) (asymmetric.IPrivKey, error) {
	keyFile := filesystem.OpenFile(pKeyPath)

	if keyFile.IsExist() {
		privKeyStr, err := keyFile.Read()
		if err != nil {
			return nil, errors.WrapError(err, "read private key")
		}
		privKey := asymmetric.LoadRSAPrivKey(string(privKeyStr))
		if privKey == nil {
			return nil, errors.NewError("private key is invalid")
		}
		return privKey, nil
	}

	privKey := asymmetric.NewRSAPrivKey(pCfg.GetSettings().GetKeySizeBits())
	if err := keyFile.Write([]byte(privKey.ToString())); err != nil {
		return nil, errors.WrapError(err, "write private key")
	}

	return privKey, nil
}
