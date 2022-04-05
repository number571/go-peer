package logger

import (
	"os"
	"strings"
	"testing"

	"github.com/number571/go-peer/cmd/hms/utils"
)

const (
	tcPathInfo    = "logger_info_test.txt"
	tcPathWarning = "logger_warning_test.txt"
	tcPathError   = "logger_error_test.txt"
)

const (
	tcTestInfo    = "test_info_text"
	tcTestWarning = "test_warning_text"
	tcTestError   = "test_error_text"
)

func TestLogger(t *testing.T) {
	defer func() {
		os.Remove(tcPathInfo)
		os.Remove(tcPathWarning)
		os.Remove(tcPathError)
	}()

	fileInfo, err := os.OpenFile(tcPathInfo, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	fileWarning, err := os.OpenFile(tcPathWarning, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	fileError, err := os.OpenFile(tcPathError, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	logger := NewLogger(fileInfo, fileWarning, fileError)

	logger.Info(tcTestInfo)
	logger.Warning(tcTestWarning)
	logger.Error(tcTestError)

	if !strings.Contains(string(utils.ReadFile(tcPathInfo)), tcTestInfo) {
		t.Errorf("info does not contains tcTestInfo")
	}

	if !strings.Contains(string(utils.ReadFile(tcPathWarning)), tcTestWarning) {
		t.Errorf("warning does not contains tcTestWarning")
	}

	if !strings.Contains(string(utils.ReadFile(tcPathError)), tcTestError) {
		t.Errorf("error does not contains tcTestError")
	}
}
