package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

var (
	_ Cipher = &CipherAES{}
)

type CipherAES struct {
	key []byte
}

func NewCipher(key []byte) Cipher {
	return &CipherAES{
		key: HashSum(key),
	}
}

func (cph *CipherAES) Encrypt(msg []byte) []byte {
	block, err := aes.NewCipher(cph.key)
	if err != nil {
		return nil
	}
	blockSize := block.BlockSize()
	msg = paddingPKCS5(msg, blockSize)
	cipherText := make([]byte, blockSize+len(msg))
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], msg)
	return cipherText
}

func (cph *CipherAES) Decrypt(msg []byte) []byte {
	block, err := aes.NewCipher(cph.key)
	if err != nil {
		return nil
	}
	blockSize := block.BlockSize()
	if len(msg) < blockSize {
		return nil
	}
	iv := msg[:blockSize]
	msg = msg[blockSize:]
	if len(msg)%blockSize != 0 {
		return nil
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(msg, msg)
	return unpaddingPKCS5(msg)
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
