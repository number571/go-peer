package handler

import "errors"

var (
	ErrReadOnlineConnections = errors.New("read online connections")
	ErrReadConnections       = errors.New("read connections")
	ErrGetAllConnections     = errors.New("get all connections")
	ErrGetPublicKey          = errors.New("get public key")
	ErrGetSettingsHLS        = errors.New("get settings hls")
)
