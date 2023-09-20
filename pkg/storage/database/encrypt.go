package database

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/errors"
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
		return nil, errors.NewError("incorrect size of encrypted data")
	}

	encDataBytes := pEncBytes[hashing.CSHA256Size:]

	gotHash := pEncBytes[:hashing.CSHA256Size]
	newHash := hashing.NewHMACSHA256Hasher(
		pAuthKey,
		encDataBytes,
	).ToBytes()

	if !bytes.Equal(gotHash, newHash) {
		return nil, errors.NewError("incorrect hash of decrypted data")
	}

	decBytes := pCipher.DecryptBytes(encDataBytes)
	if decBytes == nil {
		return nil, errors.NewError("failed decrypt message")
	}

	return decBytes, nil
}
