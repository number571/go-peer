package gopeer

import (
	"crypto/rsa"
	"net"
	"sync"
)

type Listener struct {
	address string
	client  *Client
	listen  net.Listener
}

type Client struct {
	mutex       *sync.Mutex
	mapping     map[string]bool
	publicKey   *rsa.PublicKey
	privateKey  *rsa.PrivateKey
	connections map[net.Conn]string
	actions     map[string]chan bool
	F2F         FriendToFriend
}

type FriendToFriend struct {
	Enabled bool
	friends map[string]*rsa.PublicKey
}

type Package struct {
	Info InfoPackage `json:"info"`
	Head HeadPackage `json:"head"`
	Body BodyPackage `json:"body"`
}

type InfoPackage struct {
	Network string `json:"network"`
}

type HeadPackage struct {
	Rand    string `json:"rand"`
	Title   string `json:"title"`
	Sender  string `json:"sender"`
	Session string `json:"session"`
}

type BodyPackage struct {
	Data string `json:"data"`
	Hash string `json:"hash"`
	Sign string `json:"sign"`
}
