package anonymity

import "errors"

const (
	errPrefix = "pkg/network/anonymity = "
)

var (
	ErrSetHashIntoDB          = errors.New(errPrefix + "set hash into database")
	ErrNilDB                  = errors.New(errPrefix + "database is nil")
	ErrEnqueueMessage         = errors.New(errPrefix + "enqueue message")
	ErrEncryptPayload         = errors.New(errPrefix + "encrypt payload")
	ErrUnknownType            = errors.New(errPrefix + "unknown type")
	ErrLoadMessage            = errors.New(errPrefix + "load message")
	ErrStoreHashWithBroadcast = errors.New(errPrefix + "store hash with broadcast")
	ErrActionIsNotFound       = errors.New(errPrefix + "action is not found")
	ErrActionIsClosed         = errors.New(errPrefix + "action is closed")
	ErrActionTimeout          = errors.New(errPrefix + "action timeout")
	ErrEnqueuePayload         = errors.New(errPrefix + "enqueue payload")
	ErrFetchResponse          = errors.New(errPrefix + "fetch response")
	ErrBroadcastPayload       = errors.New(errPrefix + "broadcast payload")
	ErrRunning                = errors.New(errPrefix + "node running")
)
