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
	ErrUnknownMessageType    = &SMessageError{"unknown type of message"}
	ErrLoadMessageBytes      = &SMessageError{"load message bytes"}
	ErrKeySizeGteMessageSize = &SMessageError{"key size >= message size"}
)
