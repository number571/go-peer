package storage

const (
	errPrefix = "cmd/hidden_lake/helpers/traffic/internal/storage = "
)

type SStorageError struct {
	str string
}

func (err *SStorageError) Error() string {
	return errPrefix + err.str
}

var (
	ErrMessageIsExist    = &SStorageError{"message is exist"}
	ErrMessageIsNotExist = &SStorageError{"message is not exist"}
	ErrLoadMessage       = &SStorageError{"load message"}
	ErrInvalidKeySize    = &SStorageError{"invalid key size"}
	ErrSetHashIntoDB     = &SStorageError{"set hash into db"}
	ErrGetHashFromDB     = &SStorageError{"get hash from db"}
	ErrHashAlreadyExist  = &SStorageError{"hash already exist"}
)
