package app

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/number571/go-peer/internal/flag"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/types"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

// initApp work with the raw data = read files, read args
func InitApp(pArgs []string, pDefaultPath, pDefaultKey string, pParallel uint64) (types.IRunner, error) {
	strParallel := flag.GetFlagValue(pArgs, "parallel", strconv.FormatUint(pParallel, 10))
	setParallel, err := strconv.Atoi(strParallel)
	if err != nil {
		return nil, fmt.Errorf("set parallel: %w", err)
	}
	if setParallel == 0 {
		return nil, errors.New("set parallel = 0")
	}

	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", pDefaultPath), "/")
	inputKey := flag.GetFlagValue(pArgs, "key", pDefaultKey)

	cfg, err := config.InitConfig(fmt.Sprintf("%s/%s", inputPath, pkg_settings.CPathYML), nil)
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

	return NewApp(cfg, privKey, inputPath, uint64(setParallel)), nil
}

func getPrivKey(pCfg config.IConfig, pKeyPath string) (asymmetric.IPrivKey, error) {
	if _, err := os.Stat(pKeyPath); os.IsNotExist(err) {
		privKey := asymmetric.NewRSAPrivKey(pCfg.GetSettings().GetKeySizeBits())
		if err := os.WriteFile(pKeyPath, []byte(privKey.ToString()), 0o644); err != nil {
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
