package utils

import "errors"

const (
	errPrefix = "cmd/hidden_lake/applications/filesharer/internal/utils = "
)

var (
	ErrRespSizeGeLimit = errors.New(errPrefix + "response size >= limit message size")
	ErrGetSettingsHLS  = errors.New(errPrefix + "get settings from HLS")
	ErrGetSizeInBase64 = errors.New(errPrefix + "get size in base64")
)
