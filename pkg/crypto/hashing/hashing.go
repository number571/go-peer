package hashing

import (
	"crypto/sha256"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IHasher = &sSHA256Hasher{}
)

const (
	CSHA256Size = sha256.Size
)

type sSHA256Hasher struct {
	fHash []byte
}

func NewSHA256Hasher(pData []byte) IHasher {
	h := sha256.New()
	h.Write(pData)
	return &sSHA256Hasher{
		fHash: h.Sum(nil),
	}
}

func (p *sSHA256Hasher) ToString() string {
	return encoding.HexEncode(p.ToBytes())
}

func (p *sSHA256Hasher) ToBytes() []byte {
	return p.fHash
}
