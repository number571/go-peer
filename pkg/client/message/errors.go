package message

const (
	errPrefix = "pkg/client/message = "
)

type SMessageError struct {
	str string
}

func (err *SMessageError) Error() string { return errPrefix + err.str }

var (
	ErrUnknownType        = &SMessageError{"unknown type of message"}
	ErrSeparatorNotFound  = &SMessageError{"separator is not found"}
	ErrDeserializeMessage = &SMessageError{"deserialize message"}
	ErrDecodePayload      = &SMessageError{"decode hex payload"}
	ErrInvalidMessage     = &SMessageError{"invalid message"}
)
