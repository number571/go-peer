package logger

import (
	"os"

	"github.com/number571/go-peer/pkg/logger"
)

func DefaultLogger(logging ILogging) logger.ILogger {
	return logger.NewLogger(defaultSettings(logging))
}

func defaultSettings(logging ILogging) logger.ISettings {
	sett := &logger.SSettings{}
	if logging.Info() {
		sett.FInfo = os.Stdout
	}
	if logging.Warn() {
		sett.FInfo = os.Stderr
	}
	if logging.Erro() {
		sett.FInfo = os.Stderr
	}
	return logger.NewSettings(sett)
}
