package message

import "errors"

var (
	ErrUnknownType        = errors.New("unknown type of message")
	ErrSeparatorNotFound  = errors.New("separator is not found")
	ErrDeserializeMessage = errors.New("deserialize message")
	ErrDecodePayload      = errors.New("decode hex payload")
	ErrInvalidMessage     = errors.New("invalid message")
)
