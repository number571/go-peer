package conn

const (
	errPrefix = "pkg/network/conn = "
)

type SConnError struct {
	str string
}

func (err *SConnError) Error() string {
	return errPrefix + err.str
}

var (
	ErrReadFromSocket      = &SConnError{"read from socket"}
	ErrWriteToSocket       = &SConnError{"write to socket"}
	ErrReadHeaderBlock     = &SConnError{"read header block"}
	ErrInvalidMsgSize      = &SConnError{"invalid msgSize"}
	ErrReadHeaderBytes     = &SConnError{"read header bytes"}
	ErrReadBodyBytes       = &SConnError{"read body bytes"}
	ErrInvalidMessageBytes = &SConnError{"invalid message bytes"}
	ErrSendPayloadBytes    = &SConnError{"send payload bytes"}
	ErrCreateConnection    = &SConnError{"create connection"}
	ErrSetReadDeadline     = &SConnError{"set read deadline"}
	ErrSetWriteDeadline    = &SConnError{"set write deadline"}
)
