package crypto

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/number571/go-peer/encoding"
)

var (
	_ IHasher = &sSHA256{}
	_ IHasher = &sHMAC256{}
)

const (
	TruncatedSize = 20
	HashSize      = 32
	HashKeyType   = "go-peer\\sha256"
	HmacKeyType   = "go-peer\\hmac-sha256"
)

type sSHA256 struct {
	fHash []byte
}

func NewHasher(data []byte) IHasher {
	h := sha256.New()
	h.Write(data)
	return &sSHA256{
		fHash: h.Sum(nil),
	}
}

func (h *sSHA256) String() string {
	bytes := h.Bytes()
	return encoding.Base64Encode(bytes[:TruncatedSize])
}

func (h *sSHA256) Bytes() []byte {
	return h.fHash
}

func (h *sSHA256) Type() string {
	return HashKeyType
}

func (h *sSHA256) Size() uint64 {
	return HashSize
}

type sHMAC256 struct {
	fHash []byte
}

func NewHasherMAC(key []byte, data []byte) IHasher {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return &sHMAC256{
		fHash: h.Sum(nil),
	}
}

func (h *sHMAC256) String() string {
	bytes := h.Bytes()
	return encoding.Base64Encode(bytes[:TruncatedSize])
}

func (h *sHMAC256) Bytes() []byte {
	return h.fHash
}

func (h *sHMAC256) Type() string {
	return HmacKeyType
}

func (h *sHMAC256) Size() uint64 {
	return HashSize
}

func sizeInBits(n uint) uint {
	return n * 8
}
