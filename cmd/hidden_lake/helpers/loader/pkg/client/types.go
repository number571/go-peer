package client

type IClient interface {
	GetIndex() (string, error)

	RunTransfer() error
	StopTransfer() error
}

type IRequester interface {
	GetIndex() (string, error)

	RunTransfer() error
	StopTransfer() error
}
