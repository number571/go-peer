package logger

import (
	"fmt"
	"log"
)

var (
	_ ILogger = &sLogger{}
)

type sLogger struct {
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

func NewLogger(sett ISettings) ILogger {
	logger := &sLogger{}

	infoStream := sett.GetStreamInfo()
	if infoStream != nil {
		logger.fInfoOut = log.New(infoStream, fmt.Sprintf("%s[INFO] %s", colorCyan, colorReset), log.LstdFlags)
	}

	warnStream := sett.GetStreamWarn()
	if warnStream != nil {
		logger.fWarnOut = log.New(warnStream, fmt.Sprintf("%s[WARN] %s", colorYellow, colorReset), log.LstdFlags)
	}

	erroStream := sett.GetStreamErro()
	if erroStream != nil {
		logger.fErroOut = log.New(erroStream, fmt.Sprintf("%s[ERRO] %s", colorRed, colorReset), log.LstdFlags)
	}

	return logger
}

func (l *sLogger) Info(info string) {
	if l.fInfoOut == nil {
		return
	}
	l.fInfoOut.Println(info)
}

func (l *sLogger) Warn(warn string) {
	if l.fWarnOut == nil {
		return
	}
	l.fWarnOut.Println(warn)
}

func (l *sLogger) Erro(erro string) {
	if l.fErroOut == nil {
		return
	}
	l.fErroOut.Println(erro)
}
