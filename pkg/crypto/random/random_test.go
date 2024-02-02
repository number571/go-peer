package random

import (
	"bytes"
	"testing"
)

func TestStdPRNG(t *testing.T) {
	t.Parallel()

	r := NewStdPRNG()

	if bytes.Equal(r.GetBytes(8), r.GetBytes(8)) {
		t.Error("bytes in random equals")
	}

	x := r.GetString(8)
	if x == r.GetString(8) {
		t.Error("strings in random equals")
	}

	y := r.GetUint64()
	if y == r.GetUint64() {
		t.Error("numbers in random equals")
	}
}

func TestStdPRNGBool(t *testing.T) {
	t.Parallel()

	r := NewStdPRNG()
	for i := 0; i < 1000; i++ {
		t1 := r.GetBool()
		t2 := r.GetBool()
		if t1 != t2 {
			break
		}
	}
}
