package logger

import (
	"io"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FInfo io.Writer
	FWarn io.Writer
	FErro io.Writer
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

func (p *sSettings) GetInfoWriter() io.Writer {
	return p.FInfo
}

func (p *sSettings) GetWarnWriter() io.Writer {
	return p.FWarn
}

func (p *sSettings) GetErroWriter() io.Writer {
	return p.FErro
}
