package database

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
)

func doEncrypt(pCipher symmetric.ICipher, pAuthKey []byte, pDataBytes []byte) []byte {
	encDataBytes := pCipher.EncryptBytes(pDataBytes)
	return bytes.Join(
		[][]byte{
			hashing.NewHMACSHA256Hasher(
				pAuthKey,
				encDataBytes,
			).ToBytes(),
			encDataBytes,
		},
		[]byte{},
	)
}

func tryDecrypt(pCipher symmetric.ICipher, pAuthKey []byte, pEncBytes []byte) ([]byte, error) {
	if len(pEncBytes) < hashing.CSHA256Size+symmetric.CAESBlockSize {
		return nil, ErrInvalidEncryptedSize
	}

	encDataBytes := pEncBytes[hashing.CSHA256Size:]

	gotHash := pEncBytes[:hashing.CSHA256Size]
	newHash := hashing.NewHMACSHA256Hasher(
		pAuthKey,
		encDataBytes,
	).ToBytes()

	if !bytes.Equal(gotHash, newHash) {
		return nil, ErrInvalidDataHash
	}

	return pCipher.DecryptBytes(encDataBytes), nil
}
