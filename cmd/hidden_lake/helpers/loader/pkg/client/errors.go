package client

import "errors"

var (
	ErrRequest        = errors.New("request")
	ErrDecodeResponse = errors.New("decode response")
	ErrInvalidTitle   = errors.New("invalid title")
)
