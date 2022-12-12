package hashing

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/number571/go-peer/pkg/encoding"
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

func (h *sHMACSHA256Hasher) String() string {
	return encoding.HexEncode(h.Bytes())
}

func (h *sHMACSHA256Hasher) Bytes() []byte {
	return h.fHash
}

func (h *sHMACSHA256Hasher) Type() string {
	return CHMACSHA256HmacKeyType
}

func (h *sHMACSHA256Hasher) Size() uint64 {
	return CSHA256Size
}
