package network

import "errors"

const (
	errPrefix = "pkg/network = "
)

var (
	ErrNoConnections        = errors.New(errPrefix + "no connections")
	ErrWriteTimeout         = errors.New(errPrefix + "write timeout")
	ErrBroadcastMessage     = errors.New(errPrefix + "broadcast message")
	ErrCreateListener       = errors.New(errPrefix + "create listener")
	ErrListenerAccept       = errors.New(errPrefix + "listener accept")
	ErrHasLimitConnections  = errors.New(errPrefix + "has limit connections")
	ErrConnectionIsExist    = errors.New(errPrefix + "connection already exist")
	ErrConnectionIsNotExist = errors.New(errPrefix + "connection is not exist")
	ErrCloseConnection      = errors.New(errPrefix + "close connection")
	ErrAddConnections       = errors.New(errPrefix + "add connection")
)
