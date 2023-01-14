package config

type IConfig interface {
	Address() iAddress
	Connection() string
	StorageKey() string
}

type iAddress interface {
	Interface() string
	Incoming() string
}
