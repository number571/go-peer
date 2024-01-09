package logger

import (
	"fmt"
	"log"
)

var (
	_ ILogger = &sLogger{}
)

type sLogger struct {
	fSettings ISettings

	fLogFunc ILogFunc
	fInfoOut *log.Logger
	fWarnOut *log.Logger
	fErroOut *log.Logger
}

const (
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorReset  = "\033[0m"
)

func NewLogger(pSett ISettings, pLogFunc ILogFunc) ILogger {
	logger := &sLogger{
		fSettings: pSett,
		fLogFunc:  pLogFunc,
	}

	infoStream := pSett.GetStreamInfo()
	if infoStream != nil {
		logger.fInfoOut = log.New(infoStream, fmt.Sprintf("%s[INFO] %s", colorCyan, colorReset), log.LstdFlags)
	}

	warnStream := pSett.GetStreamWarn()
	if warnStream != nil {
		logger.fWarnOut = log.New(warnStream, fmt.Sprintf("%s[WARN] %s", colorYellow, colorReset), log.LstdFlags)
	}

	erroStream := pSett.GetStreamErro()
	if erroStream != nil {
		logger.fErroOut = log.New(erroStream, fmt.Sprintf("%s[ERRO] %s", colorRed, colorReset), log.LstdFlags)
	}

	return logger
}

func (p *sLogger) GetSettings() ISettings {
	return p.fSettings
}

func (p *sLogger) PushInfo(pMsg ILogArg) {
	if p.fInfoOut == nil {
		return
	}
	p.fInfoOut.Println(p.fLogFunc(pMsg))
}

func (p *sLogger) PushWarn(pMsg ILogArg) {
	if p.fWarnOut == nil {
		return
	}
	p.fWarnOut.Println(p.fLogFunc(pMsg))
}

func (p *sLogger) PushErro(pMsg ILogArg) {
	if p.fErroOut == nil {
		return
	}
	p.fErroOut.Println(p.fLogFunc(pMsg))
}
