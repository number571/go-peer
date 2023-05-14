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

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FInfo: pSett.FInfo,
		FWarn: pSett.FWarn,
		FErro: pSett.FErro,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	// set nil for void fields
	return p
}

func (p *sSettings) GetStreamInfo() *os.File {
	return p.FInfo
}

func (p *sSettings) GetStreamWarn() *os.File {
	return p.FWarn
}

func (p *sSettings) GetStreamErro() *os.File {
	return p.FErro
}
