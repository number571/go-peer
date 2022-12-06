package symmetric

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/number571/go-peer/modules/crypto/hashing"
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
	nmsg := make([]byte, len(msg))
	copy(nmsg, msg)
	block, err := aes.NewCipher(cph.fKey)
	if err != nil {
		return nil
	}
	blockSize := block.BlockSize()
	nmsg = paddingPKCS5(nmsg, blockSize)
	cipherText := make([]byte, blockSize+len(nmsg))
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], nmsg)
	return cipherText
}

func (cph *sAESCipher) Decrypt(msg []byte) []byte {
	nmsg := make([]byte, len(msg))
	copy(nmsg, msg)
	block, err := aes.NewCipher(cph.fKey)
	if err != nil {
		return nil
	}
	blockSize := block.BlockSize()
	if len(nmsg) < blockSize {
		return nil
	}
	iv := nmsg[:blockSize]
	nmsg = nmsg[blockSize:]
	if len(nmsg)%blockSize != 0 {
		return nil
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(nmsg, nmsg)
	return unpaddingPKCS5(nmsg)
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
	return CAESBlockSize
}

func paddingPKCS5(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func unpaddingPKCS5(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return nil
	}
	unpadding := int(origData[length-1])
	if length < unpadding {
		return nil
	}
	return origData[:(length - unpadding)]
}
