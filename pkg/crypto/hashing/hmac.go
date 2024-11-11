package hashing

import (
	"crypto/hmac"
	"crypto/sha512"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IHasher = &sHMACSHA512Hasher{}
)

type sHMACSHA512Hasher struct {
	fOnce    sync.Once
	fHash    []byte
	fHashStr string
}

func NewHMACHasher(pKey []byte, pData []byte) IHasher {
	h := hmac.New(sha512.New384, pKey)
	h.Write(pData)
	s := h.Sum(nil)
	return &sHMACSHA512Hasher{fHash: s}
}

func (p *sHMACSHA512Hasher) ToString() string {
	p.fOnce.Do(func() {
		p.fHashStr = encoding.HexEncode(p.ToBytes())
	})
	return p.fHashStr
}

func (p *sHMACSHA512Hasher) ToBytes() []byte {
	return p.fHash
}
