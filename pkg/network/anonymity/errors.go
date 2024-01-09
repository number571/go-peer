package anonymity

import "errors"

var (
	ErrCloseWrapperDB         = errors.New("close wrapper database")
	ErrNetworkBroadcast       = errors.New("network broadcast")
	ErrSetHashIntoDB          = errors.New("set hash into database")
	ErrNilDB                  = errors.New("database is nil")
	ErrEnqueueMessage         = errors.New("enqueue message")
	ErrEncryptPayload         = errors.New("encrypt payload")
	ErrUnknownType            = errors.New("unknown type")
	ErrLoadMessage            = errors.New("load message")
	ErrStoreHashWithBroadcast = errors.New("store hash with broadcast")
	ErrActionIsNotFound       = errors.New("action is not found")
	ErrActionIsClosed         = errors.New("action is closed")
	ErrActionTimeout          = errors.New("action timeout")
	ErrEnqueuePayload         = errors.New("enqueue payload")
	ErrFetchResponse          = errors.New("fetch response")
	ErrBroadcastPayload       = errors.New("broadcast payload")
	ErrRunning                = errors.New("node running")
)
