package hashing

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/number571/go-peer/encoding"
)

var (
	_ IHasher = &sSHA256Hasher{}
	_ IHasher = &sHMACSHA256Hasher{}
)

const (
	HashSize    = 32
	HashKeyType = "go-peer/sha256"
	HmacKeyType = "go-peer/hmac-sha256"
)

type sSHA256Hasher struct {
	fHash []byte
}

func NewSHA256Hasher(data []byte) IHasher {
	h := sha256.New()
	h.Write(data)
	return &sSHA256Hasher{
		fHash: h.Sum(nil),
	}
}

func (h *sSHA256Hasher) String() string {
	bytes := h.Bytes()
	return encoding.Base64Encode(bytes)
}

func (h *sSHA256Hasher) Bytes() []byte {
	return h.fHash
}

func (h *sSHA256Hasher) Type() string {
	return HashKeyType
}

func (h *sSHA256Hasher) Size() uint64 {
	return HashSize
}

type sHMACSHA256Hasher struct {
	fHash []byte
}

func NewHMACSHA256Hasher(key []byte, data []byte) IHasher {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return &sHMACSHA256Hasher{
		fHash: h.Sum(nil),
	}
}

func (h *sHMACSHA256Hasher) String() string {
	bytes := h.Bytes()
	return encoding.Base64Encode(bytes)
}

func (h *sHMACSHA256Hasher) Bytes() []byte {
	return h.fHash
}

func (h *sHMACSHA256Hasher) Type() string {
	return HmacKeyType
}

func (h *sHMACSHA256Hasher) Size() uint64 {
	return HashSize
}
