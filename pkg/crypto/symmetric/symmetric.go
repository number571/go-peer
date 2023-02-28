package symmetric

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
)

var (
	_ ICipher = &sAESCipher{}
)

const (
	CAESBlockSize = 16
	CAESKeySize   = hashing.CSHA256Size
	CAESKeyType   = "go-peer/aes"
)

type sAESCipher struct {
	fKey []byte
}

func NewAESCipher(key []byte) ICipher {
	return &sAESCipher{
		fKey: hashing.NewSHA256Hasher(key).ToBytes(),
	}
}

func (cph *sAESCipher) EncryptBytes(msg []byte) []byte {
	block, err := aes.NewCipher(cph.fKey)
	if err != nil {
		return nil
	}

	blockSize := block.BlockSize()
	iv := random.NewStdPRNG().GetBytes(uint64(blockSize))

	stream := cipher.NewCTR(block, iv)
	result := make([]byte, len(msg)+len(iv))
	copy(result[:blockSize], iv)

	stream.XORKeyStream(result[blockSize:], msg)
	return result
}

func (cph *sAESCipher) DecryptBytes(msg []byte) []byte {
	block, err := aes.NewCipher(cph.fKey)
	if err != nil {
		return nil
	}

	blockSize := block.BlockSize()
	if len(msg) < blockSize {
		return nil
	}

	stream := cipher.NewCTR(block, msg[:blockSize])
	result := make([]byte, len(msg)-blockSize)

	stream.XORKeyStream(result, msg[blockSize:])
	return result
}

func (cph *sAESCipher) ToString() string {
	return fmt.Sprintf("Key(%s){%X}", cph.GetType(), cph.ToBytes())
}

func (cph *sAESCipher) ToBytes() []byte {
	return cph.fKey
}

func (cph *sAESCipher) GetType() string {
	return CAESKeyType
}

func (cph *sAESCipher) GetSize() uint64 {
	return CAESKeySize
}
