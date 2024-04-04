package language

import "errors"

const (
	errPrefix = "internal/language = "
)

var (
	ErrUnknownLanguage = errors.New(errPrefix + "unknown language")
)
