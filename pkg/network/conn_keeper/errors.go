package conn_keeper // nolint: revive

import "errors"

var (
	ErrRunning = errors.New("conn_keeper running")
)
