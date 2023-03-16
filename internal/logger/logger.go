package logger

import (
	"os"

	"github.com/number571/go-peer/pkg/logger"
)

func StdLogger(logging ILogging) logger.ILogger {
	return logger.NewLogger(defaultSettings(logging))
}

func defaultSettings(logging ILogging) logger.ISettings {
	sett := &logger.SSettings{}
	if logging.HasInfo() {
		sett.FInfo = os.Stdout
	}
	if logging.HasWarn() {
		sett.FWarn = os.Stdout
	}
	if logging.HasErro() {
		sett.FErro = os.Stderr
	}
	return logger.NewSettings(sett)
}
