package gopeer

import (
	"crypto/rsa"
	"net"
	"sync"
)

// Basic structure describing the user.
// Stores the private key and list of friends.
type Client struct {
	handle      func(*Client, *Package)
	mutex       sync.Mutex
	privateKey  *rsa.PrivateKey
	mapping     map[string]bool
	connections map[string]net.Conn
	actions     map[string]chan string
	F2F         *friendToFriend
}

type friendToFriend struct {
	mutex   sync.Mutex
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
	Npow uint64 `json:"npow"`
}
