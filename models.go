package gopeer

import (
	"crypto/rsa"
	"net"
)

/* BEGIN PACKAGE PART */
type Package struct {
	Info Info
	From From
	To   To
	Head Head
	Body Body
}

type Info struct {
	Network string
	Version string
}

type Head struct {
	Title  string
	Option string
}

type From struct {
	Sender  Sender
	Address string
}

type Sender struct {
	Hashname string
	Public   string
}

type To struct {
	Receiver Receiver
	Address  string
}

type Receiver struct {
	Hashname string
}

type Body struct {
	Data string
	Desc Desc
}

type Desc struct {
	Rand       string
	PrevHash   string
	CurrHash   string
	Sign       string
	Nonce      uint64
	Difficulty uint8
}

/* END PACKAGE PART */

/* BEGIN LISTENER PART */
type Listener struct {
	Address Address
	Setting Setting
	Clients map[string]*Client
}

type Address struct {
	Ipv4 string
	Port string
}

type Setting struct {
	Listen net.Listener
}

type Client struct {
	Hashname    string
	Keys        Keys
	Address     string
	Connections map[string]*Connect
}

type Keys struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

type Connect struct {
	Connected   bool
	Session     []byte
	PrevSession []byte
	Waiting     chan bool
	Address     string
	LastHash    string
	Public      *rsa.PublicKey
}

/* END LISTENER PART */

type Destination struct {
	Address string
	Public  *rsa.PublicKey
}
