package layer1

const (
	errPrefix = "pkg/message/layer1 = "
)

type SMessageError struct {
	str string
}

func (err *SMessageError) Error() string {
	return errPrefix + err.str
}

var (
	ErrUnknownType        = &SMessageError{"unknown type"}
	ErrInvalidHeaderSize  = &SMessageError{"length of message bytes < size of header"}
	ErrInvalidProofOfWork = &SMessageError{"got invalid proof of work"}
	ErrDecodeBytesJoiner  = &SMessageError{"decode bytes joiner"}
	ErrInvalidPayloadSize = &SMessageError{"got invalid payload size"}
	ErrInvalidAuthHash    = &SMessageError{"got invalid auth hash"}
	ErrInvalidTimestamp   = &SMessageError{"got invalid timestamp"}
	ErrDecodePayload      = &SMessageError{"decode payload"}
)
