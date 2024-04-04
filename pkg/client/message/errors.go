package message

import "errors"

const (
	errPrefix = "pkg/client/message = "
)

var (
	ErrUnknownType        = errors.New(errPrefix + "unknown type of message")
	ErrSeparatorNotFound  = errors.New(errPrefix + "separator is not found")
	ErrDeserializeMessage = errors.New(errPrefix + "deserialize message")
	ErrDecodePayload      = errors.New(errPrefix + "decode hex payload")
	ErrInvalidMessage     = errors.New(errPrefix + "invalid message")
)
