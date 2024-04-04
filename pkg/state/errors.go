package state

import "errors"

const (
	errPrefix = "pkg/state = "
)

var (
	ErrStateEnabled  = errors.New(errPrefix + "state already enabled")
	ErrStateDisabled = errors.New(errPrefix + "state already disabled")
	ErrFuncEnable    = errors.New(errPrefix + "enable state function")
	ErrFuncDisable   = errors.New(errPrefix + "disable state function")
)
