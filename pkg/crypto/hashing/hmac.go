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
	fHash    []byte
	fHashStr string
}

func NewHMACHasher(pKey []byte, pData []byte) IHasher {
	h := hmac.New(sha512.New384, pKey)
	h.Write(pData)
	s := h.Sum(nil)
	return &sHMACSHA512Hasher{
		fHash:    s,
		fHashStr: encoding.HexEncode(s),
	}
}

func (p *sHMACSHA512Hasher) ToString() string {
	return p.fHashStr
}

func (p *sHMACSHA512Hasher) ToBytes() []byte {
	return p.fHash
}
