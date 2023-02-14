package config

type IConfig interface {
	Address() IAddress
	Connection() IConnection
	StorageKey() string
}

type IAddress interface {
	Interface() string
	Incoming() string
}

type IConnection interface {
	Service() string
	Traffic() string
}
