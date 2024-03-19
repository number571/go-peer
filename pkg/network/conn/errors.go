package conn

import "errors"

var (
	ErrReadFromSocket      = errors.New("read from socket")
	ErrWriteToSocket       = errors.New("write to socket")
	ErrReadHeaderBlock     = errors.New("read header block")
	ErrInvalidMsgSize      = errors.New("invalid msgSize")
	ErrReadHeaderBytes     = errors.New("read header bytes")
	ErrReadBodyBytes       = errors.New("read body bytes")
	ErrInvalidMessageBytes = errors.New("invalid message bytes")
	ErrSendPayloadBytes    = errors.New("send payload bytes")
	ErrCreateConnection    = errors.New("create connection")
)
