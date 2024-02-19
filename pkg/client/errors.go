package client

import "errors"

var (
	ErrLimitMessageSize     = errors.New("limit message size")
	ErrInitCheckMessage     = errors.New("init check message")
	ErrDecryptCipherKey     = errors.New("decrypt cipher key")
	ErrDecryptPublicKey     = errors.New("decrypt public key")
	ErrInvalidPublicKeySize = errors.New("invalid public key size")
	ErrDecodePayloadWrapper = errors.New("decode payload wrapper")
	ErrInvalidDataHash      = errors.New("invalid data hash")
	ErrInvalidHashSign      = errors.New("invalid hash sign")
	ErrInvalidPayloadSize   = errors.New("invalid payload size")
	ErrDecodePayload        = errors.New("decode payload")
	ErrEncryptSymmetricKey  = errors.New("encrypt symmetric key")
)
