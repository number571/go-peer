package queue

import "errors"

var (
	ErrRunning    = errors.New("queue running")
	ErrQueueLimit = errors.New("queue limit")
)
