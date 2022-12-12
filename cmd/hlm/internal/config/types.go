package config

type IConfig interface {
	Address() iAddress
	Connection() string
}

type iAddress interface {
	Interface() string
	Incoming() string
}
