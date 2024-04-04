package client

import "errors"

const (
	errPrefix = "pkg/client = "
)

var (
	ErrLimitMessageSize     = errors.New(errPrefix + "limit message size")
	ErrInitCheckMessage     = errors.New(errPrefix + "init check message")
	ErrDecryptCipherKey     = errors.New(errPrefix + "decrypt cipher key")
	ErrDecryptPublicKey     = errors.New(errPrefix + "decrypt public key")
	ErrInvalidPublicKeySize = errors.New(errPrefix + "invalid public key size")
	ErrDecodePayloadWrapper = errors.New(errPrefix + "decode payload wrapper")
	ErrInvalidDataHash      = errors.New(errPrefix + "invalid data hash")
	ErrInvalidHashSign      = errors.New(errPrefix + "invalid hash sign")
	ErrInvalidPayloadSize   = errors.New(errPrefix + "invalid payload size")
	ErrDecodePayload        = errors.New(errPrefix + "decode payload")
	ErrEncryptSymmetricKey  = errors.New(errPrefix + "encrypt symmetric key")
)
