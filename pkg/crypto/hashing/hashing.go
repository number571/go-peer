package hashing

import (
	"crypto/sha512"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IHasher = &sSHA512Hasher{}
)

const (
	CHasherSize = sha512.Size
)

type sSHA512Hasher struct {
	fHash []byte
}

func NewHasher(pData []byte) IHasher {
	h := sha512.New()
	h.Write(pData)
	return &sSHA512Hasher{
		fHash: h.Sum(nil),
	}
}

func (p *sSHA512Hasher) ToString() string {
	return encoding.HexEncode(p.ToBytes())
}

func (p *sSHA512Hasher) ToBytes() []byte {
	return p.fHash
}
