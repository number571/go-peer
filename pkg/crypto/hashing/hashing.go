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

func NewSHA256Hasher(data []byte) IHasher {
	h := sha256.New()
	h.Write(data)
	return &sSHA256Hasher{
		fHash: h.Sum(nil),
	}
}

func (h *sSHA256Hasher) ToString() string {
	return fmt.Sprintf("Hash(%s){%X}", h.GetType(), h.ToBytes())
}

func (h *sSHA256Hasher) ToBytes() []byte {
	return h.fHash
}

func (h *sSHA256Hasher) GetType() string {
	return CSHA256KeyType
}

func (h *sSHA256Hasher) GetSize() uint64 {
	return CSHA256Size
}
