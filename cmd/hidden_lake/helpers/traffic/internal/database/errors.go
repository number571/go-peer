package database

const (
	errPrefix = "cmd/hidden_lake/helpers/traffic/internal/database = "
)

type SDatabaseError struct {
	str string
}

func (err *SDatabaseError) Error() string {
	return errPrefix + err.str
}

var (
	ErrMessageIsExist     = &SDatabaseError{"message is exist"}
	ErrMessageIsNotExist  = &SDatabaseError{"message is not exist"}
	ErrGtMessagesCapacity = &SDatabaseError{"gt message capacity"}
	ErrInvalidKeySize     = &SDatabaseError{"invalid key size"}
	ErrLoadMessage        = &SDatabaseError{"load message"}
	ErrCloseDB            = &SDatabaseError{"close db"}
	ErrSetPointer         = &SDatabaseError{"set pointer"}
	ErrIncrementPointer   = &SDatabaseError{"increment pointer"}
	ErrWriteMessage       = &SDatabaseError{"write message"}
	ErrRewriteKeyHash     = &SDatabaseError{"rewrite key hash"}
	ErrDeleteOldKey       = &SDatabaseError{"delete old key"}
	ErrCreateDB           = &SDatabaseError{"create db"}
)
