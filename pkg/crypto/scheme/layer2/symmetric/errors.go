package symmetric

const (
	errPrefix = "pkg/crypto/scheme/layer2/symmetric = "
)

type SClientError struct {
	str string
}

func (err *SClientError) Error() string {
	return errPrefix + err.str
}

var (
	ErrInvalidKeyType   = &SClientError{"invalid key type"}
	ErrDecryptMessage   = &SClientError{"decrypt message"}
	ErrLimitMessageSize = &SClientError{"limit message size"}
	ErrMessageSize      = &SClientError{"message size"}
	ErrDecodeMessage    = &SClientError{"decode message"}
)
