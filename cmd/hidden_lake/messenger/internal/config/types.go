package config

type IConfig interface {
	Address() iAddress
	Connection() iConnection
	StorageKey() string
}

type iConnection interface {
	Service() string
	Traffic() string
}

type iAddress interface {
	Interface() string
	Incoming() string
}
