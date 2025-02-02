package client

const (
	errPrefix = "pkg/client = "
)

type SClientError struct {
	str string
}

func (err *SClientError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLimitMessageSize     = &SClientError{"limit message size"}
	ErrInitCheckMessage     = &SClientError{"init check message"}
	ErrDecryptCipherKey     = &SClientError{"decrypt cipher key"}
	ErrDecodePublicKey      = &SClientError{"decode public key"}
	ErrDecodePayloadWrapper = &SClientError{"decode payload wrapper"}
	ErrInvalidDataHash      = &SClientError{"invalid data hash"}
	ErrInvalidHashSign      = &SClientError{"invalid hash sign"}
	ErrEncryptSymmetricKey  = &SClientError{"encrypt symmetric key"}
	ErrDecodeBytesJoiner    = &SClientError{"decode bytes joiner"}
)
