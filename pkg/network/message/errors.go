package message

const (
	errPrefix = "pkg/network/message = "
)

type SMessageError struct {
	str string
}

func (err *SMessageError) Error() string { return errPrefix + err.str }

var (
	ErrUnknownType        = &SMessageError{"unknown type"}
	ErrInvalidHeaderSize  = &SMessageError{"length of message bytes < size of header"}
	ErrInvalidProofOfWork = &SMessageError{"got invalid proof of work"}
	ErrInvalidPayloadSize = &SMessageError{"got invalid payload size"}
	ErrInvalidAuthHash    = &SMessageError{"got invalid auth hash"}
	ErrDecodePayload      = &SMessageError{"decode payload"}
)
