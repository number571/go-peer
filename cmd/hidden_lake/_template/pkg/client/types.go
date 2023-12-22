package client

type IClient interface {
	GetIndex() (string, error)
	// TODO: need implementation
}

type IRequester interface {
	GetIndex() (string, error)
	// TODO: need implementation
}
