package message

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	// first digits of PI
	authSalt = "1415926535_8979323846_2643383279_5028841971_6939937510"
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	fPayload     payload.IPayload
	fHashPayload []byte
}

func NewMessage(pPld payload.IPayload, pNetworkKey string) IMessage {
	return &sMessage{
		fPayload:     pPld,
		fHashPayload: getHash(pNetworkKey, pPld.ToBytes()),
	}
}

func LoadMessage(pData []byte, pNetworkKey string) IMessage {
	// check Hash[uN]
	if len(pData) < hashing.CSHA256Size {
		return nil
	}

	hashRecv := pData[:hashing.CSHA256Size]
	pldBytes := pData[hashing.CSHA256Size:]

	if !bytes.Equal(
		hashRecv,
		getHash(pNetworkKey, pldBytes),
	) {
		return nil
	}

	// check Head[u64]
	pld := payload.LoadPayload(pldBytes)
	if pld == nil {
		return nil
	}

	return &sMessage{
		fPayload:     pld,
		fHashPayload: hashRecv,
	}
}

func (p *sMessage) GetHash() []byte {
	return p.fHashPayload
}

func (p *sMessage) GetPayload() payload.IPayload {
	return p.fPayload
}

func (p *sMessage) ToBytes() []byte {
	return bytes.Join(
		[][]byte{
			p.fHashPayload,
			p.fPayload.ToBytes(),
		},
		[]byte{},
	)
}

func getHash(pNetworkKey string, pBytes []byte) []byte {
	authKey := keybuilder.NewKeyBuilder(1, []byte(authSalt)).Build(pNetworkKey)
	return hashing.NewHMACSHA256Hasher(authKey, pBytes).ToBytes()
}
