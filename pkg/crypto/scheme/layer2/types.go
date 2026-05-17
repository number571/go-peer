package layer2

import "github.com/number571/go-peer/pkg/types"

type (
	IKeysContainer  interface{}
	IParticipantKey types.IConverter
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
