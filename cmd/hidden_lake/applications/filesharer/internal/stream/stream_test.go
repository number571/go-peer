package stream

import "testing"

func TestError(t *testing.T) {
	str := "value"
	err := &SStreamError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestNothing(_ *testing.T) {}
