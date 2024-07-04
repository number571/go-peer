package app

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/settings"
	"github.com/number571/go-peer/internal/flag"
	"github.com/number571/go-peer/internal/initapp"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"
)

// initApp work with the raw data = read files, read args
func InitApp(pArgs []string, pDefaultPath, pDefaultKey string, pDefaultParallel uint64) (types.IRunner, error) {
	strParallel := flag.GetFlagValue(pArgs, "parallel", strconv.FormatUint(pDefaultParallel, 10))
	setParallel, err := strconv.Atoi(strParallel)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetParallelValue, err)
	}
	if setParallel == 0 {
		return nil, ErrSetParallelNull
	}

	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", pDefaultPath), "/")
	cfg, err := config.InitConfig(filepath.Join(inputPath, settings.CPathYML), nil)
	if err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}

	inputKey := flag.GetFlagValue(pArgs, "key", pDefaultKey)
	privKey, err := initapp.GetPrivKey(inputKey, cfg.GetSettings().GetKeySizeBits())
	if err != nil {
		return nil, utils.MergeErrors(ErrGetPrivateKey, err)
	}

	return NewApp(cfg, privKey, uint64(setParallel)), nil
}
