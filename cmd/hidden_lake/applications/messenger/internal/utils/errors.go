package utils

import "errors"

var (
	ErrGetSizeInBase64    = errors.New("get size in base64")
	ErrMessageSizeGtLimit = errors.New("message size > limit")
	ErrGetSettingsHLS     = errors.New("get settings hls")
)
