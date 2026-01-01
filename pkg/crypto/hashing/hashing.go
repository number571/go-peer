package hashing

import (
	"crypto/sha512"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IHasher = &sSHA512Hasher{}
)

const (
	CHasherSize = sha512.Size384
)

type sSHA512Hasher struct {
	fOnce    sync.Once
	fHash    []byte
	fHashStr string
}

func NewHasher(pData interface{}) IHasher {
	var d []byte
	switch x := pData.(type) {
	case []byte:
		d = x
	case string:
		d = []byte(x)
	default:
		panic("invalid type of data")
	}
	s := sha512.Sum384(d)
	return &sSHA512Hasher{fHash: s[:]}
}

func (p *sSHA512Hasher) ToString() string {
	p.fOnce.Do(func() {
		p.fHashStr = encoding.HexEncode(p.ToBytes())
	})
	return p.fHashStr
}

func (p *sSHA512Hasher) ToBytes() []byte {
	return p.fHash
}
