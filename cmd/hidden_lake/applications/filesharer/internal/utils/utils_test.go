package utils

import "testing"

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SUtilsError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestNothing(_ *testing.T) {}
