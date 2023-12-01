package message

import (
	"github.com/number571/go-peer/pkg/types"
)

type IMessage interface {
	types.IConverter

	IsValid(ISettings) bool
	GetPubKey() []byte  // Public key of the sender.
	GetEncKey() []byte  // One-time key of encryption data.
	GetSalt() []byte    // Random bytes for hide data of the hash.
	GetHash() []byte    // Hash of the (sender + receiver + payload).
	GetSign() []byte    // Sign of the hash.
	GetPayload() []byte // Main data in the ecnrypted bytes format.
}

type ISettings interface {
	GetKeySizeBits() uint64
	GetMessageSizeBytes() uint64
}
