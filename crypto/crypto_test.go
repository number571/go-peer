package crypto

import (
	"bytes"
	"testing"
)

func TestHash(t *testing.T) {
	msg := []byte("hello, world!")

	hash := NewHasher(msg).String()
	if hash != NewHasher(msg).String() {
		t.Errorf("hash is not determined")
	}

	msg[3] = msg[3] ^ 8
	if hash == NewHasher(msg).String() {
		t.Errorf("bit didn't change the result ")
	}
}

func TestSign(t *testing.T) {
	var (
		priv = NewPrivKey(1024)
		msg  = []byte("hello, world!")
	)

	pub := priv.PubKey()
	sign := priv.Sign(msg)

	if !pub.Verify(msg, sign) {
		t.Errorf("signature is invalid")
	}
}

func TestAEncrypt(t *testing.T) {
	var (
		priv = NewPrivKey(1024)
		msg  = []byte("hello, world!")
	)

	pub := priv.PubKey()
	emsg := pub.Encrypt(msg)

	if !bytes.Equal(msg, priv.Decrypt(emsg)) {
		t.Errorf("decrypted message is invalid")
	}
}

func TestSEncrypt(t *testing.T) {
	var (
		key = []byte("it's a key!")
		msg = []byte("hello, world!")
	)

	cipher := NewCipher(key)
	emsg := cipher.Encrypt(msg)

	if !bytes.Equal(msg, cipher.Decrypt(emsg)) {
		t.Errorf("decrypted message is invalid")
	}
}

func TestPuzzle(t *testing.T) {
	var (
		puzzle = NewPuzzle(10)
		msg    = []byte("hello, world!")
	)

	hash := NewHasher(msg).Bytes()
	proof := puzzle.Proof(hash)

	if !puzzle.Verify(hash, proof) {
		t.Errorf("proof is invalid")
	}

	hash[3] = hash[3] ^ 8
	if puzzle.Verify(hash, proof) {
		t.Errorf("proof is correct?")
	}
}

func TestRand(t *testing.T) {
	if bytes.Equal(RandBytes(8), RandBytes(8)) {
		t.Errorf("bytes in random equals")
	}

	if RandString(8) == RandString(8) {
		t.Errorf("strings in random equals")
	}

	if RandUint64() == RandUint64() {
		t.Errorf("numbers in random equals")
	}
}

func TestEntropy(t *testing.T) {
	var (
		msg  = []byte("hello, world!")
		salt = []byte("it's a salt!")
	)

	hash := RaiseEntropy(msg, salt, 10)

	if bytes.Equal(hash, NewHasher(msg).Bytes()) {
		t.Errorf("hash is correct?")
	}

	if !bytes.Equal(hash, RaiseEntropy(msg, salt, 10)) {
		t.Errorf("hash is not determined")
	}
}
