package mempool

import "github.com/number571/go-peer/cmd/union_blockchain/kernel/block"

var (
	_ ISettings = &sSettings{}
)

const (
	cPath     = "mempool.db"
	cCountTXs = 512
)

type SSettings sSettings
type sSettings struct {
	FCountTXs      uint64
	FBlockSettings block.ISettings
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FCountTXs:      sett.FCountTXs,
		FBlockSettings: sett.FBlockSettings,
	}).useDefaultValues()
}

func (s *sSettings) useDefaultValues() ISettings {
	if s.FCountTXs == 0 {
		s.FCountTXs = cCountTXs
	}
	if s.FBlockSettings == nil {
		s.FBlockSettings = block.NewSettings(&block.SSettings{})
	}
	return s
}

func (sett *sSettings) GetCountTXs() uint64 {
	return sett.FCountTXs
}

func (sett *sSettings) GetBlockSettings() block.ISettings {
	return sett.FBlockSettings
}
