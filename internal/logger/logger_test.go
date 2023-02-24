package logger

import (
	"os"
	"testing"
)

var (
	_ = &tsLogger{}
)

type tsLogger struct{}

func (l *tsLogger) Info() bool {
	return true
}

func (l *tsLogger) Warn() bool {
	return true
}

func (l *tsLogger) Erro() bool {
	return false
}

func TestLogger(t *testing.T) {
	logger := DefaultLogger(&tsLogger{})
	if logger.Settings().GetStreamInfo().Name() != os.Stdout.Name() {
		t.Error("info stream != stdout")
		return
	}
	if logger.Settings().GetStreamWarn().Name() != os.Stderr.Name() {
		t.Error("warn stream != stderr")
		return
	}
	if logger.Settings().GetStreamErro() != nil {
		t.Error("erro stream != nil")
		return
	}
}
