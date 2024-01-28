package client

import "errors"

var (
	ErrRequest        = errors.New("request")
	ErrDecodeResponse = errors.New("decode response")
)
