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
func InitApp(pArgs []string, pDefaultPath string, pDefaultParallel uint64) (types.IRunner, error) {
	strParallel := flag.GetFlagValue(pArgs, "parallel", strconv.FormatUint(pDefaultParallel, 10))
	setParallel, err := strconv.ParseUint(strParallel, 10, 64)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetParallel, err)
	}

	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", pDefaultPath), "/")

	cfgPath := filepath.Join(inputPath, pkg_settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil)
	if err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}

	keyPath := filepath.Join(inputPath, pkg_settings.CPathKey)
	privKey, err := initapp.GetPrivKey(keyPath, cfg.GetSettings().GetKeySizeBits())
	if err != nil {
		return nil, utils.MergeErrors(ErrGetPrivateKey, err)
	}

	return NewApp(cfg, privKey, inputPath, setParallel), nil
}
