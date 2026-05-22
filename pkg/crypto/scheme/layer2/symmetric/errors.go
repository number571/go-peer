package symmetric

const (
	errPrefix = "pkg/crypto/scheme/layer2/symmetric = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrInvalidKeyType     = &SError{"invalid key type"}
	ErrDecryptMessage     = &SError{"decrypt message"}
	ErrLimitMessageSize   = &SError{"limit message size"}
	ErrInvalidMessageSize = &SError{"invalidmessage size"}
	ErrDecodeMessage      = &SError{"decode message"}
)
