package client

import "errors"

var (
	ErrBadRequest       = errors.New("bad request")
	ErrDecodeResponse   = errors.New("decode response")
	ErrInvalidPublicKey = errors.New("invalid public key")
	ErrInvalidTitle     = errors.New("invalid title")
)
