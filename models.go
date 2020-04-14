package gopeer

import (
	"crypto/rsa"
	"crypto/x509"
	"net"
	"sync"
)

type conndata struct {
	Certificate string
	Public      string
	Session     string
}

type optionType uint8

const (
	_raw     optionType = 0
	_confirm optionType = 1
)

/* BEGIN PACKAGE PART */
type Package struct {
	Info Info `json:"info"`
	From From `json:"from"`
	To   To   `json:"to"`
	Head Head `json:"head"`
	Body Body `json:"body"`
}

type Info struct {
	Network string `json:"network"`
	Version string `json:"version"`
}

type Head struct {
	Title  string `json:"title"`
	Option string `json:"option"`
}

type From struct {
	Sender   Sender `json:"sender"`
	Hashname string `json:"hashname"`
	Address  string `json:"address"`
}

type To struct {
	Receiver Receiver `json:"receiver"`
	Hashname string   `json:"hashname"`
	Address  string   `json:"address"`
}

type Body struct {
	Data string `json:"data"`
	Desc Desc   `json:"desc"`
	Test Test   `json:"test"`
}

type Receiver Hidden
type Sender Hidden
type Hidden struct {
	Hashname string `json:"hashname"`
}

type Desc struct {
	Id          uint64 `json:"id"`
	Rand        string `json:"rand"`
	Hash        string `json:"hash"`
	Sign        string `json:"sign"`
	Nonce       uint64 `json:"nonce"`
	Difficulty  uint8  `json:"difficulty"`
	Redirection uint8  `json:"redirection"`
}

type Test struct {
	Hash string `json:"hash"`
	Sign string `json:"sign"`
}
/* END PACKAGE PART */

/* BEGIN LISTENER PART */
type Listener struct {
	listen      net.Listener
	handleFunc  func(*Client, *Package)
	address     address
	certificate []byte
	Clients     map[string]*Client
}

type address struct {
	ipv4 string
	port string
}

type Client struct {
	listener    *Listener
	remember    remember
	keys        keys
	hashname    string
	address     string
	certPool    *x509.CertPool
	F2F         F2F
	Sharing     Sharing
	Mutex       *sync.Mutex
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

type keys struct {
	private *rsa.PrivateKey
	public  *rsa.PublicKey
}

type Connect struct {
	connected   bool
	hashname    string
	packageId   uint64
	relation    net.Conn
	transfer    transfer
	address     string
	session     []byte
	certificate []byte
	throwClient *rsa.PublicKey
	public      *rsa.PublicKey
	Chans       Chans
}

type Chans struct {
	Action chan bool
	action chan bool
}

type transfer struct {
	active     bool
	packdata   string
}
/* END LISTENER PART */

/* BEGIN FILE TRANSFER */
type FileTransfer struct {
	Head HeadTransfer `json:"head"`
	Body BodyTransfer `json:"body"`
}

type HeadTransfer struct {
	Id     uint64 `json:"id"`
	Name   string `json:"name"`
	IsNull bool   `json:"is_null"`
}

type BodyTransfer struct {
	Hash []byte `json:"hash"`
	Data []byte `json:"data"`
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
