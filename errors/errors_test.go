package errors

import (
	"fmt"
	"testing"
)

const (
	ctErrorMsg = "example of custom error"
)

func TestErrors(t *testing.T) {
	err := NewError(1, ctErrorMsg)
	if err == nil {
		t.Error("error is not nil")
		return
	}

	if !IsError(err, 1) {
		t.Error("type of error != 1")
		return
	}

	if err.Error() != ctErrorMsg {
		t.Error("message of error not equal")
		return
	}

	err = fmt.Errorf("new standart error")
	if IsError(err, 1) {
		t.Error("type of error = 1 but is a standart error")
		return
	}
}
