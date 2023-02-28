package hashing

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

var (
	_ IHasher = &sHMACSHA256Hasher{}
)

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

func (h *sHMACSHA256Hasher) ToString() string {
	return fmt.Sprintf("HMAC(%s){%X}", h.GetType(), h.ToBytes())
}

func (h *sHMACSHA256Hasher) ToBytes() []byte {
	return h.fHash
}

func (h *sHMACSHA256Hasher) GetType() string {
	return CSHA256KeyType
}

func (h *sHMACSHA256Hasher) GetSize() uint64 {
	return CSHA256Size
}
