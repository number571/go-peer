package handler

import "errors"

const (
	errPrefix = "cmd/hidden_lake/applications/filesharer/internal/handler = "
)

var (
	ErrReadOnlineConnections = errors.New(errPrefix + "read online connections")
	ErrReadConnections       = errors.New(errPrefix + "read connections")
	ErrGetAllConnections     = errors.New(errPrefix + "get all connections")
	ErrGetPublicKey          = errors.New(errPrefix + "get public key")
	ErrGetSettingsHLS        = errors.New(errPrefix + "get settings hls")
)
