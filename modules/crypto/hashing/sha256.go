package hashing

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/number571/go-peer/modules/encoding"
	"github.com/number571/go-peer/settings"
)

var (
	_ IHasher = &sSHA256Hasher{}
	_ IHasher = &sHMACSHA256Hasher{}
)

const (
	CSHA256Size            = 32
	CSHA256KeyType         = settings.CGopeerPrefix + "sha256"
	CHMACSHA256HmacKeyType = settings.CGopeerPrefix + "hmac-sha256"
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
	return fmt.Sprintf("Hash(%s){%X}", h.Type(), h.Bytes())
}

func (h *sSHA256Hasher) Bytes() []byte {
	return h.fHash
}

func (h *sSHA256Hasher) Type() string {
	return CSHA256KeyType
}

func (h *sSHA256Hasher) Size() uint64 {
	return CSHA256Size
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
