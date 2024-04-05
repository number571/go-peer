package adapted

import "testing"

func TestError(t *testing.T) {
	str := "value"
	err := &SAdaptedError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestNothing(_ *testing.T) {}
