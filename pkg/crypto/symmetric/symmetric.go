package symmetric

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/random"
)

var (
	_ ICipher = &sAESCipher{}
)

const (
	CAESBlockSize = aes.BlockSize
	CAESKeySize   = 32
	CAESKeyType   = "go-peer/aes"
)

type sAESCipher struct {
	fKey []byte
}

func NewAESCipher(pKey []byte) ICipher {
	if len(pKey) != CAESKeySize {
		panic("len(pKey) != CAESKeySize")
	}
	return &sAESCipher{
		fKey: pKey,
	}
}

func (p *sAESCipher) EncryptBytes(pMsg []byte) []byte {
	block, err := aes.NewCipher(p.fKey)
	if err != nil {
		return nil
	}

	blockSize := block.BlockSize()
	iv := random.NewStdPRNG().GetBytes(uint64(blockSize))

	stream := cipher.NewCFBEncrypter(block, iv)
	result := make([]byte, len(pMsg)+len(iv))
	copy(result[:blockSize], iv)

	stream.XORKeyStream(result[blockSize:], pMsg)
	return result
}

func (p *sAESCipher) DecryptBytes(pMsg []byte) []byte {
	block, err := aes.NewCipher(p.fKey)
	if err != nil {
		return nil
	}

	blockSize := block.BlockSize()
	if len(pMsg) < blockSize {
		return nil
	}

	stream := cipher.NewCFBDecrypter(block, pMsg[:blockSize])
	result := make([]byte, len(pMsg)-blockSize)

	stream.XORKeyStream(result, pMsg[blockSize:])
	return result
}

func (p *sAESCipher) ToString() string {
	return fmt.Sprintf("Key(%s){%X}", p.GetType(), p.ToBytes())
}

func (p *sAESCipher) ToBytes() []byte {
	return p.fKey
}

func (p *sAESCipher) GetType() string {
	return CAESKeyType
}

func (p *sAESCipher) GetSize() uint64 {
	return CAESKeySize
}
