package hashing

import "testing"

func TestSHA256(t *testing.T) {
	t.Parallel()

	msg := []byte("hello, world!")

	hash := NewSHA256Hasher(msg).ToString()
	if hash != NewSHA256Hasher(msg).ToString() {
		t.Error("hash is not determined")
		return
	}

	msg[3] = msg[3] ^ 8
	if hash == NewSHA256Hasher(msg).ToString() {
		t.Error("bit didn't change the result ")
		return
	}
}

func TestHMACSHA256(t *testing.T) {
	t.Parallel()

	key := []byte("secret key")
	msg := []byte("hello, world!")

	hash := NewHMACSHA256Hasher(key, msg).ToString()
	if hash != NewHMACSHA256Hasher(key, msg).ToString() {
		t.Error("hash is not determined")
		return
	}

	msg[3] = msg[3] ^ 8
	if hash == NewHMACSHA256Hasher(key, msg).ToString() {
		t.Error("bit didn't change the result")
		return
	}
}
