package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	_ ILogger = &sLogger{}
)

type sLogger struct {
	fInfoOut    *log.Logger
	fWarningOut *log.Logger
	fErrorOut   *log.Logger
}

const (
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorReset  = "\033[0m"
)

func NewLogger(infoOut, warningOut, errorOut *os.File) ILogger {
	return &sLogger{
		log.New(infoOut, fmt.Sprintf("%s%s%s", colorCyan, "[INFO]\t", colorReset), log.LstdFlags),
		log.New(warningOut, fmt.Sprintf("%s%s%s", colorYellow, "[WARN]\t", colorReset), log.LstdFlags),
		log.New(errorOut, fmt.Sprintf("%s%s%s", colorRed, "[ERRO]\t", colorReset), log.LstdFlags),
	}
}

func (l *sLogger) Info(info string) {
	l.fInfoOut.Println(info)
}

func (l *sLogger) Warning(info string) {
	l.fWarningOut.Println(info)
}

func (l *sLogger) Error(info string) {
	l.fErrorOut.Println(info)
}
