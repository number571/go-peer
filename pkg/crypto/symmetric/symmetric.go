package symmetric

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/number571/go-peer/pkg/crypto/random"
)

var (
	_ ICipher = &sAESCipher{}
)

const (
	CAESBlockSize = aes.BlockSize
	CAESKeySize   = 32
)

type sAESCipher struct {
	fBlock cipher.Block
}

func NewAESCipher(pKey []byte) ICipher {
	if len(pKey) != CAESKeySize {
		panic("len(pKey) != CAESKeySize")
	}
	block, _ := aes.NewCipher(pKey)
	return &sAESCipher{
		fBlock: block,
	}
}

func (p *sAESCipher) EncryptBytes(pMsg []byte) []byte {
	blockSize := p.fBlock.BlockSize()
	iv := random.NewCSPRNG().GetBytes(uint64(blockSize))

	stream := cipher.NewCFBEncrypter(p.fBlock, iv)
	result := make([]byte, len(pMsg)+len(iv))
	copy(result[:blockSize], iv)

	stream.XORKeyStream(result[blockSize:], pMsg)
	return result
}

func (p *sAESCipher) DecryptBytes(pMsg []byte) []byte {
	blockSize := p.fBlock.BlockSize()
	if len(pMsg) < blockSize {
		return nil
	}

	stream := cipher.NewCFBDecrypter(p.fBlock, pMsg[:blockSize])
	result := make([]byte, len(pMsg)-blockSize)

	stream.XORKeyStream(result, pMsg[blockSize:])
	return result
}
