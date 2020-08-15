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
	privateKey  *rsa.PrivateKey
	connections map[net.Conn]string
	actions     map[string]chan string
	f2f         friendToFriend
}

type friendToFriend struct {
	enabled bool
	friends map[string]*rsa.PublicKey
}

type Package struct {
	Head HeadPackage `json:"head"`
	Body BodyPackage `json:"body"`
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
