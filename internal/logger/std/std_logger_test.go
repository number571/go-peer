package logger

import (
	"os"
	"testing"
)

var (
	_ = &tsLogger{}
)

type tsLogger struct{}

func (l *tsLogger) HasInfo() bool {
	return true
}

func (l *tsLogger) HasWarn() bool {
	return true
}

func (l *tsLogger) HasErro() bool {
	return false
}

func TestLogger(t *testing.T) {
	logger := StdLogger(&tsLogger{})
	if logger.GetSettings().GetStreamInfo().Name() != os.Stdout.Name() {
		t.Error("info stream != stdout")
		return
	}
	if logger.GetSettings().GetStreamWarn().Name() != os.Stdout.Name() {
		t.Error("warn stream != stdout")
		return
	}
	if logger.GetSettings().GetStreamErro() != nil {
		t.Error("erro stream != nil")
		return
	}
}
