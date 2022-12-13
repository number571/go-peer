package logger

import (
	"os"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FInfo *os.File
	FWarn *os.File
	FErro *os.File
}

func NewSettings(sett *SSettings) ISettings {
	return (&sSettings{
		FInfo: sett.FInfo,
		FWarn: sett.FWarn,
		FErro: sett.FErro,
	}).useDefaultValues()
}

func (s *sSettings) useDefaultValues() ISettings {
	// set nil for void fields
	return s
}

func (s *sSettings) GetStreamInfo() *os.File {
	return s.FInfo
}

func (s *sSettings) GetStreamWarn() *os.File {
	return s.FWarn
}

func (s *sSettings) GetStreamErro() *os.File {
	return s.FErro
}
