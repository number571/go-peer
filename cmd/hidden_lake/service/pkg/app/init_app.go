package app

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/number571/go-peer/internal/flag"
	"github.com/number571/go-peer/internal/initapp"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

// initApp work with the raw data = read files, read args
func InitApp(pArgs []string, pDefaultPath, pDefaultPrivPath string, pDefaultParallel uint64) (types.IRunner, error) {
	strParallel := flag.GetFlagValue(pArgs, "parallel", strconv.FormatUint(pDefaultParallel, 10))
	setParallel, err := strconv.Atoi(strParallel)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetParallel, err)
	}
	if setParallel == 0 {
		return nil, ErrSetParallelNull
	}

	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", pDefaultPath), "/")
	cfg, err := config.InitConfig(filepath.Join(inputPath, pkg_settings.CPathYML), nil)
	if err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}

	privPath := flag.GetFlagValue(pArgs, "priv", pDefaultPrivPath)
	privKey, err := initapp.GetPrivKey(privPath, cfg.GetSettings().GetKeySizeBits())
	if err != nil {
		return nil, utils.MergeErrors(ErrGetPrivateKey, err)
	}

	return NewApp(cfg, privKey, inputPath, uint64(setParallel)), nil
}
