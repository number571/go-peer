package queue

import "errors"

const (
	errPrefix = "pkg/network/anonymity/queue = "
)

var (
	ErrRunning    = errors.New(errPrefix + "queue running")
	ErrQueueLimit = errors.New(errPrefix + "queue limit")
)
