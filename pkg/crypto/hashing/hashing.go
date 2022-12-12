package hashing

import (
	"crypto/sha256"
	"fmt"
)

var (
	_ IHasher = &sSHA256Hasher{}
)

const (
	CSHA256Size            = 32
	CSHA256KeyType         = "go-peer/sha256"
	CHMACSHA256HmacKeyType = "go-peer/hmac-sha256"
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

func (h *sSHA256Hasher) String() string {
	return fmt.Sprintf("Hash(%s){%X}", h.Type(), h.Bytes())
}

func (h *sSHA256Hasher) Bytes() []byte {
	return h.fHash
}

func (h *sSHA256Hasher) Type() string {
	return CSHA256KeyType
}

func (h *sSHA256Hasher) Size() uint64 {
	return CSHA256Size
}
