package utils

import "errors"

const (
	errPrefix = "cmd/hidden_lake/applications/messenger/internal/utils = "
)

var (
	ErrGetSizeInBase64    = errors.New(errPrefix + "get size in base64")
	ErrMessageSizeGtLimit = errors.New(errPrefix + "message size > limit")
	ErrGetSettingsHLS     = errors.New(errPrefix + "get settings hls")
)
