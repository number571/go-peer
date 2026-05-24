package symmetric

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
)

func NewCipherCFB(pKey []byte) ICipher {
	if len(pKey) != CCipherKeySize {
		return nil
	}
	block, _ := aes.NewCipher(pKey)
	return &sAESCipher{
		fMode:   modeCFB,
		fKey:    pKey,
		fBlock:  block,
		fHasher: hashing.NewHMACHasher(pKey, []byte("__hasher__")),
	}
}

func (p *sAESCipher) encryptBytesCFB(pMsg []byte) []byte {
	blockSize := p.fBlock.BlockSize()
	iv := random.NewRandom().GetBytes(uint64(blockSize)) //nolint:gosec

	stream := cipher.NewCFBEncrypter(p.fBlock, iv)
	result := make([]byte, len(pMsg)+len(iv))
	copy(result[:blockSize], iv)

	stream.XORKeyStream(result[blockSize:], pMsg)
	return result
}

func (p *sAESCipher) decryptBytesCFB(pMsg []byte) []byte {
	blockSize := p.fBlock.BlockSize()
	if len(pMsg) < blockSize {
		return nil
	}

	stream := cipher.NewCFBDecrypter(p.fBlock, pMsg[:blockSize])
	result := make([]byte, len(pMsg)-blockSize)

	stream.XORKeyStream(result, pMsg[blockSize:])
	return result
}
