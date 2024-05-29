package handler

import "testing"

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SHandlerError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}
