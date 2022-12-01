package chain

import (
	"fmt"

	"github.com/number571/go-peer/cmd/ubc/kernel/mempool"
)

var (
	_ ISettings = &sSettings{}
)

const (
	cRootPath         = "chain.db"
	cBlocksPath       = "blocks.db"
	cTransactionsPath = "transactions.db"
	cMempoolPath      = "mempool.db"
)

type SSettings sSettings
type sSettings struct {
	FRootPath         string
	FBlocksPath       string
	FTransactionsPath string
	FMempoolPath      string
	FMempoolSettings  mempool.ISettings
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FRootPath:         sett.FRootPath,
		FBlocksPath:       sett.FBlocksPath,
		FTransactionsPath: sett.FTransactionsPath,
		FMempoolPath:      sett.FMempoolPath,
		FMempoolSettings:  sett.FMempoolSettings,
	}).useDefaultValues()
}

func (s *sSettings) useDefaultValues() ISettings {
	if s.FRootPath == "" {
		s.FRootPath = cRootPath
	}
	if s.FBlocksPath == "" {
		s.FBlocksPath = cBlocksPath
	}
	if s.FTransactionsPath == "" {
		s.FTransactionsPath = cTransactionsPath
	}
	if s.FMempoolPath == "" {
		s.FMempoolPath = cMempoolPath
	}
	if s.FMempoolSettings == nil {
		s.FMempoolSettings = mempool.NewSettings(&mempool.SSettings{})
	}
	return s
}

func (sett *sSettings) GetRootPath() string {
	return sett.FRootPath
}

func (sett *sSettings) GetBlocksPath() string {
	return fmt.Sprintf("%s/%s", sett.FRootPath, sett.FBlocksPath)
}

func (sett *sSettings) GetTransactionsPath() string {
	return fmt.Sprintf("%s/%s", sett.FRootPath, sett.FTransactionsPath)
}

func (sett *sSettings) GetMempoolPath() string {
	return fmt.Sprintf("%s/%s", sett.FRootPath, sett.FMempoolPath)
}

func (sett *sSettings) GetMempoolSettings() mempool.ISettings {
	return sett.FMempoolSettings
}
