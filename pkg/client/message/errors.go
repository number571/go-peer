package message

const (
	errPrefix = "pkg/client/message = "
)

type SMessageError struct {
	str string
}

func (err *SMessageError) Error() string {
	return errPrefix + err.str
}

var (
	ErrUnknownType     = &SMessageError{"unknown type of message"}
	ErrInvalidMessage  = &SMessageError{"invalid message"}
	ErrLoadBytesJoiner = &SMessageError{"load bytes joiner"}
)
