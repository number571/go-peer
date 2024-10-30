package hashing

import (
	"crypto/sha512"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IHasher = &sSHA512Hasher{}
)

const (
	CHasherSize = sha512.Size384
)

type sSHA512Hasher struct {
	fHash    []byte
	fHashStr string
}

func NewHasher(pData []byte) IHasher {
	s := sha512.Sum384(pData)
	return &sSHA512Hasher{
		fHash:    s[:],
		fHashStr: encoding.HexEncode(s[:]),
	}
}

func (p *sSHA512Hasher) ToString() string {
	return p.fHashStr
}

func (p *sSHA512Hasher) ToBytes() []byte {
	return p.fHash
}
