package network

import "errors"

var (
	ErrNoConnections        = errors.New("no connections")
	ErrWriteTimeout         = errors.New("write timeout")
	ErrBroadcastMessage     = errors.New("broadcast message")
	ErrCreateListener       = errors.New("create listener")
	ErrListenerAccept       = errors.New("listener accept")
	ErrHasLimitConnections  = errors.New("has limit connections")
	ErrConnectionIsExist    = errors.New("connection already exist")
	ErrConnectionIsNotExist = errors.New("connection is not exist")
	ErrCloseConnection      = errors.New("close connection")
	ErrAddConnections       = errors.New("add connection")
)
