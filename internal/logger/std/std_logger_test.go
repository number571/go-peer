package std

import (
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
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
	return true
}

func TestGetLogFunc(t *testing.T) {
	t.Parallel()

	f := GetLogFunc()
	if f("string") != "string" {
		t.Error("incorrect logger work")
		return
	}
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	_ = f(struct{}{})
}

func TestLogger(t *testing.T) {
	t.Parallel()

	logger := NewStdLogger(
		&tsLogger{},
		func(_ logger.ILogArg) string {
			return ""
		},
	)
	if logger.GetSettings().GetStreamInfo().Name() != os.Stdout.Name() {
		t.Error("info stream != stdout")
		return
	}
	if logger.GetSettings().GetStreamWarn().Name() != os.Stdout.Name() {
		t.Error("warn stream != stdout")
		return
	}
	if logger.GetSettings().GetStreamErro().Name() != os.Stderr.Name() {
		t.Error("erro stream != stderr")
		return
	}
}
