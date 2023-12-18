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

func NewHMACSHA256Hasher(pKey []byte, pData []byte) IHasher {
	h := hmac.New(sha256.New, pKey)
	h.Write(pData)
	return &sHMACSHA256Hasher{
		fHash: h.Sum(nil),
	}
}

func (p *sHMACSHA256Hasher) ToString() string {
	return encoding.HexEncode(p.ToBytes())
}

func (p *sHMACSHA256Hasher) ToBytes() []byte {
	return p.fHash
}
