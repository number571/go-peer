package crypto

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/number571/go-peer/encoding"
)

var (
	_ Hasher = &sha256T{}
	_ Hasher = &hmac256T{}
)

const (
	TruncatedSize = 20
	HashSize      = 32
	HashKeyType   = "go-peer\\sha256"
	HmacKeyType   = "go-peer\\hmac-sha256"
)

type sha256T struct {
	hash []byte
}

func NewHasher(data []byte) Hasher {
	h := sha256.New()
	h.Write(data)
	return &sha256T{
		hash: h.Sum(nil),
	}
}

func (h *sha256T) String() string {
	bytes := h.Bytes()
	return encoding.Base64Encode(bytes[:TruncatedSize])
}

func (h *sha256T) Bytes() []byte {
	return h.hash
}

func (h *sha256T) Type() string {
	return HashKeyType
}

func (h *sha256T) Size() uint64 {
	return HashSize
}

type hmac256T struct {
	hash []byte
}

func NewHasherMAC(key []byte, data []byte) Hasher {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return &hmac256T{
		hash: h.Sum(nil),
	}
}

func (h *hmac256T) String() string {
	bytes := h.Bytes()
	return encoding.Base64Encode(bytes[:TruncatedSize])
}

func (h *hmac256T) Bytes() []byte {
	return h.hash
}

func (h *hmac256T) Type() string {
	return HmacKeyType
}

func (h *hmac256T) Size() uint64 {
	return HashSize
}

func sizeInBits(n uint) uint {
	return n * 8
}
