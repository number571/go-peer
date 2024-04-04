package message

import "errors"

const (
	errPrefix = "pkg/network/message = "
)

var (
	ErrUnknownType        = errors.New(errPrefix + "unknown type")
	ErrInvalidHeaderSize  = errors.New(errPrefix + "length of message bytes < size of header")
	ErrInvalidProofOfWork = errors.New(errPrefix + "got invalid proof of work")
	ErrInvalidPayloadSize = errors.New(errPrefix + "got invalid payload size")
	ErrInvalidAuthHash    = errors.New(errPrefix + "got invalid auth hash")
	ErrDecodePayload      = errors.New(errPrefix + "decode payload")
)
