package symmetric

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/number571/go-peer/modules/crypto/hashing"
	"github.com/number571/go-peer/modules/crypto/random"
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
		fKey: hashing.NewSHA256Hasher(key).Bytes(),
	}
}

func (cph *sAESCipher) Encrypt(msg []byte) []byte {
	block, err := aes.NewCipher(cph.fKey)
	if err != nil {
		return nil
	}

	blockSize := block.BlockSize()
	iv := random.NewStdPRNG().Bytes(uint64(blockSize))

	stream := cipher.NewCTR(block, iv)
	result := make([]byte, len(msg)+len(iv))
	copy(result[:blockSize], iv)

	stream.XORKeyStream(result[blockSize:], msg)
	return result
}

func (cph *sAESCipher) Decrypt(msg []byte) []byte {
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

func (cph *sAESCipher) String() string {
	return fmt.Sprintf("Key(%s){%X}", cph.Type(), cph.Bytes())
}

func (cph *sAESCipher) Bytes() []byte {
	return cph.fKey
}

func (cph *sAESCipher) Type() string {
	return CAESKeyType
}

func (cph *sAESCipher) Size() uint64 {
	return CAESKeySize
}
