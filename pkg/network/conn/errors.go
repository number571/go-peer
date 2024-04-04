package conn

import "errors"

const (
	errPrefix = "pkg/network/conn = "
)

var (
	ErrReadFromSocket      = errors.New(errPrefix + "read from socket")
	ErrWriteToSocket       = errors.New(errPrefix + "write to socket")
	ErrReadHeaderBlock     = errors.New(errPrefix + "read header block")
	ErrInvalidMsgSize      = errors.New(errPrefix + "invalid msgSize")
	ErrReadHeaderBytes     = errors.New(errPrefix + "read header bytes")
	ErrReadBodyBytes       = errors.New(errPrefix + "read body bytes")
	ErrInvalidMessageBytes = errors.New(errPrefix + "invalid message bytes")
	ErrSendPayloadBytes    = errors.New(errPrefix + "send payload bytes")
	ErrCreateConnection    = errors.New(errPrefix + "create connection")
)
