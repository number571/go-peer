package hashing

import (
	"testing"
)

func TestHasher(t *testing.T) {
	t.Parallel()

	msg := "hello, world!"
	msgBytes := []byte(msg)

	if NewHasher(msg).ToString() != NewHasher(msgBytes).ToString() {
		t.Fatal("hash invalid with same data")
	}

	hash := NewHasher(msg).ToString()
	if hash != NewHasher(msg).ToString() {
		t.Fatal("hash is not determinated")
	}

	msgBytes[3] ^= 8
	if hash == NewHasher(msgBytes).ToString() {
		t.Fatal("bit didn't change the result")
	}
}

func TestHMACSHasher(t *testing.T) {
	t.Parallel()

	key := []byte("secret key")
	msg := []byte("hello, world!")

	hash := NewHMACHasher(key, msg).ToString()
	if hash != NewHMACHasher(key, msg).ToString() {
		t.Error("hash is not determined")
		return
	}

	msg[3] ^= 8
	if hash == NewHMACHasher(key, msg).ToString() {
		t.Error("bit didn't change the result")
		return
	}
}
