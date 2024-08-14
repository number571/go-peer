package client

type IClient interface {
	GetSettings() ISettings
	GetMessageLimit() uint64

	MessageIsValid([]byte) bool
	EncryptMessage([]byte, []byte) ([]byte, error)
	DecryptMessage([]byte, []byte) ([]byte, error)
}

type ISettings interface {
	GetMessageSizeBytes() uint64
}
