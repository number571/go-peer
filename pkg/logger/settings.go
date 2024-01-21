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

func (p *sSettings) GetOutInfo() io.Writer {
	return p.FInfo
}

func (p *sSettings) GetOutWarn() io.Writer {
	return p.FWarn
}

func (p *sSettings) GetOutErro() io.Writer {
	return p.FErro
}
