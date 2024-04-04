package client

import "errors"

var (
	ErrRequest          = errors.New("request")
	ErrDecodeResponse   = errors.New("decode response")
	ErrInvalidResponse  = errors.New("invalid response")
	ErrInvalidHexFormat = errors.New("invalid hex format")
	ErrInvalidTitle     = errors.New("invalid title")
	ErrDecodeMessage    = errors.New("decode message")
)
