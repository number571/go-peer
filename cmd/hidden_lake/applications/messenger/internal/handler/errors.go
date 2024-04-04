package handler

import "errors"

var (
	ErrReadConnections       = errors.New("read connections")
	ErrReadOnlineConnections = errors.New("read online connections")
	ErrGetAllConnections     = errors.New("get all connections")
	ErrGetSettings           = errors.New("get settings")
	ErrGetPublicKey          = errors.New("get public key")
	ErrUnknownMessageType    = errors.New("unknown message type")
	ErrUnwrapFile            = errors.New("unwrap file")
	ErrHasNotWritableChars   = errors.New("had not writable chars")
	ErrMessageNull           = errors.New("message null")
	ErrUndefinedPublicKey    = errors.New("undefined public key")
	ErrGetFriends            = errors.New("get friends")
	ErrLenMessageGtLimit     = errors.New("len message > limit")
	ErrGetMessageLimit       = errors.New("get message limit")
	ErrPushMessage           = errors.New("push message")
	ErrReadFile              = errors.New("read file")
	ErrReadFileSize          = errors.New("read file size")
	ErrGetFormFile           = errors.New("get form file")
	ErrUploadFile            = errors.New("upload file")
)
