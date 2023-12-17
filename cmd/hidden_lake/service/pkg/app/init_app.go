package app

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/number571/go-peer/internal/flag"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/types"

	pkg_config "github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

// initApp work with the raw data = read files, read args
func InitApp(pDefaultPath, pDefaultKey string) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(flag.GetFlagValue("path", pDefaultPath), "/")
	inputKey := flag.GetFlagValue("key", pDefaultKey)

	cfg, err := pkg_config.InitConfig(fmt.Sprintf("%s/%s", inputPath, pkg_settings.CPathYML), nil)
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	privKey, err := getPrivKey(cfg, inputKey)
	if err != nil {
		return nil, fmt.Errorf("get private key: %w", err)
	}

	if privKey.GetSize() != cfg.GetSettings().GetKeySizeBits() {
		return nil, errors.New("size of private key is invalid")
	}

	return NewApp(cfg, privKey, inputPath), nil
}

func getPrivKey(pCfg pkg_config.IConfig, pKeyPath string) (asymmetric.IPrivKey, error) {
	if _, err := os.Stat(pKeyPath); os.IsNotExist(err) {
		privKey := asymmetric.NewRSAPrivKey(pCfg.GetSettings().GetKeySizeBits())
		if err := os.WriteFile(pKeyPath, []byte(privKey.ToString()), 0o600); err != nil {
			return nil, fmt.Errorf("write private key: %w", err)
		}
		return privKey, nil
	}

	privKeyStr, err := os.ReadFile(pKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}
	privKey := asymmetric.LoadRSAPrivKey(string(privKeyStr))
	if privKey == nil {
		return nil, errors.New("private key is invalid")
	}
	return privKey, nil
}
