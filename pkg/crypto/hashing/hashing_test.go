package hashing

import (
	"fmt"
	"testing"
)

func TestHasher(t *testing.T) {
	t.Parallel()

	msg := []byte("hello, world!")

	hash := NewHasher(msg).ToString()
	if hash != NewHasher(msg).ToString() {
		t.Error("hash is not determined")
		return
	}

	msg[3] ^= 8
	if hash == NewHasher(msg).ToString() {
		t.Error("bit didn't change the result ")
		return
	}

	fmt.Println(hash)
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
