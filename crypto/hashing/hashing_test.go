package hashing

import "testing"

func TestSHA256(t *testing.T) {
	msg := []byte("hello, world!")

	hash := NewSHA256Hasher(msg).String()
	if hash != NewSHA256Hasher(msg).String() {
		t.Errorf("hash is not determined")
	}

	msg[3] = msg[3] ^ 8
	if hash == NewSHA256Hasher(msg).String() {
		t.Errorf("bit didn't change the result ")
	}
}
