package logger

import (
	"os"
	"strings"
	"testing"

	"github.com/number571/go-peer/utils"
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
		t.Error(err.Error())
		return
	}

	fileWarning, err := os.OpenFile(tcPathWarning, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Error(err.Error())
		return
	}

	fileError, err := os.OpenFile(tcPathError, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Error(err.Error())
		return
	}

	logger := NewLogger(fileInfo, fileWarning, fileError)

	logger.Info(tcTestInfo)
	logger.Warning(tcTestWarning)
	logger.Error(tcTestError)

	res, err := utils.OpenFile(tcPathInfo).Read()
	if err != nil {
		t.Error(err.Error())
	}
	if !strings.Contains(string(res), tcTestInfo) {
		t.Error("info does not contains tcTestInfo")
	}

	res, err = utils.OpenFile(tcPathWarning).Read()
	if err != nil {
		t.Error(err.Error())
	}
	if !strings.Contains(string(res), tcTestWarning) {
		t.Error("warning does not contains tcTestWarning")
	}

	res, err = utils.OpenFile(tcPathError).Read()
	if err != nil {
		t.Error(err.Error())
	}
	if !strings.Contains(string(res), tcTestError) {
		t.Error("error does not contains tcTestError")
	}
}
