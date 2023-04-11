package message

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/payload"
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	fPayload     payload.IPayload
	fHashPayload []byte
	fEncrPayload []byte
}

func NewMessage(pPld payload.IPayload, pKey []byte) IMessage {
	cipher := symmetric.NewAESCipher(pKey)
	encPayload := cipher.EncryptBytes(pPld.ToBytes())
	return &sMessage{
		fPayload:     pPld,
		fHashPayload: getHash(cipher, encPayload),
		fEncrPayload: encPayload,
	}
}

func LoadMessage(pData, pKey []byte) IMessage {
	// check Hash[uN]
	if len(pData) < hashing.CSHA256Size {
		return nil
	}

	hashRecv := pData[:hashing.CSHA256Size]
	encPayload := pData[hashing.CSHA256Size:]

	cipher := symmetric.NewAESCipher(pKey)
	if !bytes.Equal(
		hashRecv,
		getHash(cipher, encPayload),
	) {
		return nil
	}

	// check Head[u64]
	decPayload := cipher.DecryptBytes(encPayload)
	pld := payload.LoadPayload(decPayload)
	if pld == nil {
		return nil
	}

	return &sMessage{
		fPayload:     pld,
		fHashPayload: hashRecv,
		fEncrPayload: encPayload,
	}
}

func (p *sMessage) GetHash() []byte {
	return p.fHashPayload
}

func (p *sMessage) GetPayload() payload.IPayload {
	return p.fPayload
}

func (p *sMessage) GetBytes() []byte {
	return bytes.Join(
		[][]byte{
			p.fHashPayload,
			p.fEncrPayload,
		},
		[]byte{},
	)
}

func getHash(cipher symmetric.ICipher, pBytes []byte) []byte {
	rawKey := bytes.Join(
		[][]byte{[]byte("__"), cipher.ToBytes(), []byte("__")},
		[]byte{},
	)
	return hashing.NewHMACSHA256Hasher(
		hashing.NewSHA256Hasher(rawKey).ToBytes(),
		pBytes,
	).ToBytes()
}
