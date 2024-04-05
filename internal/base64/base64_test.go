package base64

import (
	"testing"
)

func TestError(t *testing.T) {
	str := "value"
	err := &SBase64Error{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestGetSizeInBase64(t *testing.T) {
	if _, err := GetSizeInBase64(1); err == nil {
		t.Error("success get size with < 2 bytes")
		return
	}
	if n, err := GetSizeInBase64(1000); err != nil || n != 748 {
		t.Error("got invalid size in base64")
		return
	}
}
