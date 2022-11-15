package config

type IConfig interface {
	Address() iAddress
	Connection() string
}

type iAddress interface {
	WebLocal() string
	Incoming() string
}
