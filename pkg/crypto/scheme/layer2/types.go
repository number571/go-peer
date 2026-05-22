package layer2

import (
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/types"
)

type IScheme interface {
	IEncryptor
	IDecryptor

	GetRandomKey() IParticipantKey
	GetMessageSize() uint64
	GetPayloadLimit() uint64
}

type IDecryptor interface {
	DecryptMessage(IKeysContainer, []byte) (IParticipantKey, []byte, error)
}

type IEncryptor interface {
	EncryptMessage(IParticipantKey, []byte) ([]byte, error)
}

type IKeysContainer interface {
	List() []IParticipantKey
	Add(IParticipantKey) bool
	Del(string) bool
	Get(string) (IParticipantKey, bool)
}

type IParticipantKey interface {
	types.IConverter
	GetHasher() hashing.IHasher
}
