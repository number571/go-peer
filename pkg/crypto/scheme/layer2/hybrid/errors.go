package hybrid

const (
	errPrefix = "pkg/crypto/scheme/layer2/hybrid = "
)

type SSchemeError struct {
	str string
}

func (err *SSchemeError) Error() string {
	return errPrefix + err.str
}

var (
	ErrInvalidKeyType       = &SSchemeError{"invalid key type"}
	ErrLimitMessageSize     = &SSchemeError{"limit message size"}
	ErrInitCheckMessage     = &SSchemeError{"init check message"}
	ErrDecryptCipherKey     = &SSchemeError{"decrypt cipher key"}
	ErrDecodePublicKey      = &SSchemeError{"decode public key"}
	ErrDecodePayloadWrapper = &SSchemeError{"decode payload wrapper"}
	ErrInvalidDataHash      = &SSchemeError{"invalid data hash"}
	ErrInvalidHashSign      = &SSchemeError{"invalid hash sign"}
	ErrEncryptSymmetricKey  = &SSchemeError{"encrypt symmetric key"}
	ErrDecodeBytesJoiner    = &SSchemeError{"decode bytes joiner"}
	ErrUnknownMessageType   = &SSchemeError{"unknown type of message"}
	ErrLoadMessageBytes     = &SSchemeError{"load message bytes"}
	ErrSizeMessageBytes     = &SSchemeError{"size message bytes"}
	ErrStructGTEMessageSize = &SSchemeError{"struct >= message size"}
)
