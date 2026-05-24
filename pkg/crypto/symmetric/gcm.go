package symmetric

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
)

func NewCipherGCM(pKey []byte) ICipher {
	if len(pKey) != CCipherKeySize {
		return nil
	}
	block, _ := aes.NewCipher(pKey)
	return &sAESCipher{
		fMode:   modeGCM,
		fKey:    pKey,
		fBlock:  block,
		fHasher: hashing.NewHMACHasher(pKey, []byte("__hasher__")),
	}
}

func (p *sAESCipher) encryptBytesGCM(pMsg []byte) []byte {
	gcm, err := cipher.NewGCM(p.fBlock)
	if err != nil {
		return nil
	}
	nonce := random.NewRandom().GetBytes(uint64(gcm.NonceSize())) //nolint:gosec
	ciphertext := gcm.Seal(nil, nonce, pMsg, nil)
	return bytes.Join([][]byte{nonce, ciphertext}, []byte{})
}

func (p *sAESCipher) decryptBytesGCM(pMsg []byte) []byte {
	gcm, err := cipher.NewGCM(p.fBlock)
	if err != nil {
		return nil
	}
	nonceSize := gcm.NonceSize()
	if len(pMsg) < nonceSize {
		return nil
	}
	nonce, ciphertext := pMsg[:nonceSize], pMsg[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil
	}
	return plaintext
}
