package app

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/number571/go-peer/internal/flag"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

// initApp work with the raw data = read files, read args
func InitApp(pArgs []string, pDefaultPath, pDefaultKey string, pDefaultParallel uint64) (types.IRunner, error) {
	strParallel := flag.GetFlagValue(pArgs, "parallel", strconv.FormatUint(pDefaultParallel, 10))
	setParallel, err := strconv.Atoi(strParallel)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetParallel, err)
	}
	if setParallel == 0 {
		return nil, ErrSetParallelNull
	}

	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", pDefaultPath), "/")
	inputKey := flag.GetFlagValue(pArgs, "key", pDefaultKey)

	cfg, err := config.InitConfig(filepath.Join(inputPath, pkg_settings.CPathYML), nil)
	if err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}

	privKey, err := getPrivKey(cfg, inputKey)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetPrivateKey, err)
	}

	if privKey.GetSize() != cfg.GetSettings().GetKeySizeBits() {
		return nil, ErrSizePrivateKey
	}

	return NewApp(cfg, privKey, inputPath, uint64(setParallel)), nil
}

func getPrivKey(pCfg config.IConfig, pKeyPath string) (asymmetric.IPrivKey, error) {
	if _, err := os.Stat(pKeyPath); os.IsNotExist(err) {
		privKey := asymmetric.NewRSAPrivKey(pCfg.GetSettings().GetKeySizeBits())
		if err := os.WriteFile(pKeyPath, []byte(privKey.ToString()), 0o600); err != nil {
			return nil, utils.MergeErrors(ErrWritePrivateKey, err)
		}
		return privKey, nil
	}

	privKeyStr, err := os.ReadFile(pKeyPath)
	if err != nil {
		return nil, utils.MergeErrors(ErrReadPrivateKey, err)
	}
	privKey := asymmetric.LoadRSAPrivKey(string(privKeyStr))
	if privKey == nil {
		return nil, ErrInvalidPrivateKey
	}
	return privKey, nil
}
