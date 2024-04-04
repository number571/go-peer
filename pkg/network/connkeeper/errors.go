package connkeeper

import "errors"

const (
	errPrefix = "pkg/network/connkeeper = "
)

var (
	ErrRunning = errors.New(errPrefix + "connkeeper running")
)
