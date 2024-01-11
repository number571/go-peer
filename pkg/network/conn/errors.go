package conn

import "errors"

var (
	ErrReadFromSocket        = errors.New("read from socket")
	ErrWriteToSocket         = errors.New("write to socket")
	ErrReadHeaderBlock       = errors.New("read header block")
	ErrInvalidHeaderMsgSize  = errors.New("invalid header.encMsgSize")
	ErrInvalidHeaderVoidSize = errors.New("invalid header.voidSize")
	ErrInvalidHeaderAuthHash = errors.New("invalid header.authHash")
	ErrReadHeaderBytes       = errors.New("read header bytes")
	ErrReadBodyBytes         = errors.New("read body bytes")
	ErrInvalidBodyAuthHash   = errors.New("invalid body.authHash")
	ErrInvalidMessageBytes   = errors.New("invalid message bytes")
	ErrSendPayloadBytes      = errors.New("send payload bytes")
	ErrCreateConnection      = errors.New("create connection")
)
