package crypto

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/number571/go-peer/encoding"
)

var (
	_ Hasher = &SHA256{}
	_ Hasher = &HMAC256{}
)

const (
	TruncatedSize = 20
	HashSize      = 32
	HashKeyType   = "go-peer\\sha256"
	HmacKeyType   = "go-peer\\hmac-sha256"
)

type SHA256 struct {
	hash []byte
}

func NewSHA256(data []byte) Hasher {
	h := sha256.New()
	h.Write(data)
	return &SHA256{
		hash: h.Sum(nil),
	}
}

func (h *SHA256) String() string {
	bytes := h.Bytes()
	return encoding.Base64Encode(bytes[:TruncatedSize])
}

func (h *SHA256) Bytes() []byte {
	return h.hash
}

func (h *SHA256) Type() string {
	return HashKeyType
}

func (h *SHA256) Size() uint {
	return HashSize
}

type HMAC256 struct {
	hash []byte
}

func NewHMAC256(data []byte, key []byte) Hasher {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return &HMAC256{
		hash: h.Sum(nil),
	}
}

func (h *HMAC256) String() string {
	bytes := h.Bytes()
	return encoding.Base64Encode(bytes[:TruncatedSize])
}

func (h *HMAC256) Bytes() []byte {
	return h.hash
}

func (h *HMAC256) Type() string {
	return HmacKeyType
}

func (h *HMAC256) Size() uint {
	return HashSize
}
