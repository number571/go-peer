package message

import "github.com/number571/go-peer/modules/payload"

type IMessage interface {
	Hash() []byte
	Payload() payload.IPayload
	ToBytes() []byte
}

type IPackage interface {
	Size() uint64
	Data() []byte
	ToBytes() []byte
}
