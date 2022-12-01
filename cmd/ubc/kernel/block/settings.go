package block

import "github.com/number571/go-peer/cmd/ubc/kernel/transaction"

var (
	_ ISettings = &sSettings{}
)

const (
	cCountTXs = (1 << 5)
)

type SSettings sSettings
type sSettings struct {
	FCountTXs            uint64
	FTransactionSettings transaction.ISettings
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FCountTXs:            sett.FCountTXs,
		FTransactionSettings: sett.FTransactionSettings,
	}).useDefaultValues()
}

func (s *sSettings) useDefaultValues() ISettings {
	if s.FCountTXs == 0 {
		s.FCountTXs = cCountTXs
	}
	if s.FTransactionSettings == nil {
		s.FTransactionSettings = transaction.NewSettings(
			&transaction.SSettings{},
		)
	}
	return s
}

func (sett *sSettings) GetCountTXs() uint64 {
	return sett.FCountTXs
}

func (sett *sSettings) GetTransactionSettings() transaction.ISettings {
	return sett.FTransactionSettings
}
