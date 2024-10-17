package hashing

import (
	"crypto/hmac"
	"crypto/sha512"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IHasher = &sHMACSHA512Hasher{}
)

type sHMACSHA512Hasher struct {
	fHash []byte
}

func NewHMACHasher(pKey []byte, pData []byte) IHasher {
	h := hmac.New(sha512.New, pKey)
	h.Write(pData)
	return &sHMACSHA512Hasher{
		fHash: h.Sum(nil),
	}
}

func (p *sHMACSHA512Hasher) ToString() string {
	return encoding.HexEncode(p.ToBytes())
}

func (p *sHMACSHA512Hasher) ToBytes() []byte {
	return p.fHash
}
