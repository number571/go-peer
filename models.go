package gopeer

import (
	"crypto/rsa"
	"crypto/x509"
	"net"
	"sync"
)

type Option uint8
const (
	RAW Option = 0
	CONFIRM Option = 1
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

type Body struct {
	Data string
	Desc Desc
}

type Receiver Hidden
type Sender Hidden
type Hidden struct {
	Hashname string
}

type Desc struct {
	Id          uint64
	Rand        string
	Hash        string
	Sign        string
	Nonce       uint64
	Difficulty  uint8
	Redirection uint8
}
/* END PACKAGE PART */

/* BEGIN LISTENER PART */
type Listener struct {
	listen      net.Listener
	handleFunc  func(*Client, *Package)
	Address     Address
	Clients     map[string]*Client
	Certificate []byte
}

type Address struct {
	Ipv4 string
	Port string
}

type Client struct {
	listener    *Listener
	remember    remember
	F2F         F2F
	Sharing     Sharing
	Hashname    string
	Address     string
	Mutex       *sync.Mutex
	CertPool    *x509.CertPool
	Keys        Keys
	Connections map[string]*Connect
}

type remember struct {
	index   uint16
	mapping map[string]uint16
	listing []string
}

type F2F struct {
	Perm    bool
	Friends map[string]bool
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
	packageId   uint64
	relation    net.Conn
	transfer    transfer
	Chans       Chans
	Address     string
	Session     []byte
	Certificate []byte
	ThrowClient *rsa.PublicKey
	Public      *rsa.PublicKey
}

type Chans struct {
	Action chan bool
	action chan bool
}

type transfer struct {
	active bool
	inputFile  string
	outputFile string
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
	Address     string
	Certificate []byte
	Public      *rsa.PublicKey
	Receiver    *rsa.PublicKey
}

type Certificate struct {
	Cert []byte
	Key  []byte
}
