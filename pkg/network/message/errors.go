package message

import "errors"

var (
	ErrUnknownType        = errors.New("unknown type")
	ErrInvalidHeaderSize  = errors.New("length of message bytes < size of header")
	ErrInvalidProofOfWork = errors.New("got invalid proof of work")
	ErrInvalidAuthHash    = errors.New("got invalid auth hash")
	ErrDecodePayload      = errors.New("decode payload")
)
