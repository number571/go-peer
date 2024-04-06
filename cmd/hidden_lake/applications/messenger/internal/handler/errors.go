package handler

const (
	errPrefix = "cmd/hidden_lake/applications/messenger/internal/handler = "
)

type SHandlerError struct {
	str string
}

func (err *SHandlerError) Error() string {
	return errPrefix + err.str
}

var (
	ErrReadConnections       = &SHandlerError{"read connections"}
	ErrReadOnlineConnections = &SHandlerError{"read online connections"}
	ErrGetAllConnections     = &SHandlerError{"get all connections"}
	ErrGetSettings           = &SHandlerError{"get settings"}
	ErrGetPublicKey          = &SHandlerError{"get public key"}
	ErrUnknownMessageType    = &SHandlerError{"unknown message type"}
	ErrUnwrapFile            = &SHandlerError{"unwrap file"}
	ErrHasNotWritableChars   = &SHandlerError{"had not writable chars"}
	ErrMessageNull           = &SHandlerError{"message null"}
	ErrUndefinedPublicKey    = &SHandlerError{"undefined public key"}
	ErrGetFriends            = &SHandlerError{"get friends"}
	ErrLenMessageGtLimit     = &SHandlerError{"len message > limit"}
	ErrGetMessageLimit       = &SHandlerError{"get message limit"}
	ErrPushMessage           = &SHandlerError{"push message"}
	ErrReadFile              = &SHandlerError{"read file"}
	ErrReadFileSize          = &SHandlerError{"read file size"}
	ErrGetFormFile           = &SHandlerError{"get form file"}
	ErrUploadFile            = &SHandlerError{"upload file"}
)
