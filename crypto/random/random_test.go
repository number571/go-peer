package random

import (
	"bytes"
	"testing"
)

func TestStdPRNG(t *testing.T) {
	r := NewStdPRNG()

	if bytes.Equal(r.Bytes(8), r.Bytes(8)) {
		t.Errorf("bytes in random equals")
	}

	//lint:ignore SA4000 is random strings
	if r.String(8) == r.String(8) {
		t.Errorf("strings in random equals")
	}

	//lint:ignore SA4000 is random numbers
	if r.Uint64() == r.Uint64() {
		t.Errorf("numbers in random equals")
	}
}
