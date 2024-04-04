package handler

import "errors"

const (
	errPrefix = "cmd/hidden_lake/applications/messenger/internal/handler = "
)

var (
	ErrReadConnections       = errors.New(errPrefix + "read connections")
	ErrReadOnlineConnections = errors.New(errPrefix + "read online connections")
	ErrGetAllConnections     = errors.New(errPrefix + "get all connections")
	ErrGetSettings           = errors.New(errPrefix + "get settings")
	ErrGetPublicKey          = errors.New(errPrefix + "get public key")
	ErrUnknownMessageType    = errors.New(errPrefix + "unknown message type")
	ErrUnwrapFile            = errors.New(errPrefix + "unwrap file")
	ErrHasNotWritableChars   = errors.New(errPrefix + "had not writable chars")
	ErrMessageNull           = errors.New(errPrefix + "message null")
	ErrUndefinedPublicKey    = errors.New(errPrefix + "undefined public key")
	ErrGetFriends            = errors.New(errPrefix + "get friends")
	ErrLenMessageGtLimit     = errors.New(errPrefix + "len message > limit")
	ErrGetMessageLimit       = errors.New(errPrefix + "get message limit")
	ErrPushMessage           = errors.New(errPrefix + "push message")
	ErrReadFile              = errors.New(errPrefix + "read file")
	ErrReadFileSize          = errors.New(errPrefix + "read file size")
	ErrGetFormFile           = errors.New(errPrefix + "get form file")
	ErrUploadFile            = errors.New(errPrefix + "upload file")
)
