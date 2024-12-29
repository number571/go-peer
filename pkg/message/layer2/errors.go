package layer2

const (
	errPrefix = "pkg/message/layer2 = "
)

type SMessageError struct {
	str string
}

func (err *SMessageError) Error() string {
	return errPrefix + err.str
}

var (
	ErrUnknownMessageType = &SMessageError{"unknown type of message"}
	ErrLoadMessageBytes   = &SMessageError{"load message bytes"}
	ErrSizeMessageBytes   = &SMessageError{"size message bytes"}
)
