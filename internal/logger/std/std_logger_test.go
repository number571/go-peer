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
	if l := f("string"); l != "string" {
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

func TestStdLogger(t *testing.T) {
	stdLogger := NewStdLogger(&tsLogger{}, func(_ logger.ILogArg) string {
		return ""
	})
	if _, ok := stdLogger.GetSettings().GetOutInfo().(*os.File); !ok {
		t.Error("invalid info type")
		return
	}
	if _, ok := stdLogger.GetSettings().GetOutWarn().(*os.File); !ok {
		t.Error("invalid warn type")
		return
	}
	if _, ok := stdLogger.GetSettings().GetOutErro().(*os.File); !ok {
		t.Error("invalid erro type")
		return
	}
}
