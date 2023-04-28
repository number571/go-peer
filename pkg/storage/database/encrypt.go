package database

import (
	"bytes"
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
)

func doEncrypt(pCipher symmetric.ICipher, pDataBytes []byte) []byte {
	return bytes.Join(
		[][]byte{
			hashing.NewHMACSHA256Hasher(
				pCipher.ToBytes(),
				pDataBytes,
			).ToBytes(),
			pCipher.EncryptBytes(pDataBytes),
		},
		[]byte{},
	)
}

func tryDecrypt(pCipher symmetric.ICipher, pEncBytes []byte) ([]byte, error) {
	if len(pEncBytes) < hashing.CSHA256Size+symmetric.CAESBlockSize {
		return nil, fmt.Errorf("incorrect size of encrypted data")
	}

	decBytes := pCipher.DecryptBytes(pEncBytes[hashing.CSHA256Size:])
	if decBytes == nil {
		return nil, fmt.Errorf("failed decrypt message")
	}

	gotHashed := pEncBytes[:hashing.CSHA256Size]
	newHashed := hashing.NewHMACSHA256Hasher(
		pCipher.ToBytes(),
		decBytes,
	).ToBytes()

	if !bytes.Equal(gotHashed, newHashed) {
		return nil, fmt.Errorf("incorrect hash of decrypted data")
	}

	return decBytes, nil
}
