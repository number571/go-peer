package hashing

import (
	"crypto/sha256"
	"fmt"
)

var (
	_ IHasher = &sSHA256Hasher{}
)

const (
	CSHA256Size    = sha256.Size
	CSHA256KeyType = "go-peer/sha256"
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
	return fmt.Sprintf("Hash(%s){%X}", p.GetType(), p.ToBytes())
}

func (p *sSHA256Hasher) ToBytes() []byte {
	return p.fHash
}

func (p *sSHA256Hasher) GetType() string {
	return CSHA256KeyType
}

func (p *sSHA256Hasher) GetSize() uint64 {
	return CSHA256Size
}
