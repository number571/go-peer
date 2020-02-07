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
	Sender   Sender
	Hashname string
	Address  string
}

type To struct {
	Receiver Receiver
	Hashname string
	Address  string
}

type Sender Hidden
type Receiver Hidden
type Hidden struct {
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
	listen  net.Listener
	Address Address
	Clients map[string]*Client
}

type Address struct {
	Ipv4 string
	Port string
}

type Client struct {
	Hashname    string
	Address     string
	Sharing     Sharing
	Keys        Keys
	Connections map[string]*Connect
}

type Sharing struct {
	Perm bool
	Path string
}

type Keys struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

type Connect struct {
	connected   bool
	prevSession []byte
	waiting     chan bool
	lastHash    string
	transfer    transfer
	Session     []byte
	Address     string
	Public      *rsa.PublicKey
	PublicRecv  *rsa.PublicKey
}

type transfer struct {
	outputFile string
	isBlocked  bool
}

/* END LISTENER PART */

/* BEGIN FILE TRANSFER */
type FileTransfer struct {
	Head HeadTransfer
	Body BodyTransfer
}

type HeadTransfer struct {
	Id     uint32
	Name   string
	IsNull bool
}

type BodyTransfer struct {
	Hash []byte
	Data []byte
}

/* END FILE TRANSFER */

type Destination struct {
	Address  string
	Public   *rsa.PublicKey
	Receiver *rsa.PublicKey
}
