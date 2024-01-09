package state

import "errors"

var (
	ErrStateEnabled  = errors.New("state already enabled")
	ErrStateDisabled = errors.New("state already disabled")
	ErrFuncEnable    = errors.New("enable state function")
	ErrFuncDisable   = errors.New("disable state function")
)
