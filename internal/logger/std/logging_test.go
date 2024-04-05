package std

import "testing"

func TestError(t *testing.T) {
	str := "value"
	err := &SStdError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestLogging(t *testing.T) {
	logging, err := LoadLogging([]string{"info", "erro"})
	if err != nil {
		t.Error(err)
		return
	}
	if !logging.HasInfo() {
		t.Error("failed has info")
		return
	}
	if logging.HasWarn() {
		t.Error("failed has warn")
		return
	}
	if !logging.HasErro() {
		t.Error("failed has erro")
		return
	}
	if _, err := LoadLogging([]string{"info", "unknown"}); err == nil {
		t.Error("success load invalid logging")
		return
	}
}
