package gopeer

import (
	"crypto/rsa"
	"net"
	"sync"
)

type SettingsType map[string]interface{}

type RelationType uint8

const (
	RelationAll    RelationType = 0
	RelationNode   RelationType = 1
	RelationHandle RelationType = 2
	RelationHidden RelationType = 3
)

type ReadonlyType uint8

const (
	ReadAll    ReadonlyType = 0
	ReadNode   ReadonlyType = 1
	ReadHandle ReadonlyType = 2
)

type AccessType uint8

const (
	AccessDenied  AccessType = 0
	AccessAllowed AccessType = 1
)

type Node struct {
	Hashname string
	Keys     Keys
	Setting  Setting
	Address  Address
	Network  Network
}

type Setting struct {
	Mutex           *sync.Mutex
	ReadOnly        ReadonlyType
	Listen          net.Listener
	HandleServer    func(*Node, *Package)
	TestConnections map[string]bool
}

type Keys struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

type Network struct {
	Addresses   map[string]string
	AccessList  map[string]AccessType
	Connections map[string]*Connect
}

type Connect struct {
	Relation RelationType
	Hashname string
	Session  []byte
	Public   *rsa.PublicKey
	Link     net.Conn
}

type Address struct {
	IPv4 string
	Port string
}

type Package struct {
	Info Info
	From From
	To   To
	Head Head
	Body Body
}

type Info struct {
	NET string
}

type From struct {
	Hashname string
	Address  string
	Public   string
}

type To struct {
	Address string
}

type Head struct {
	Title string
	Mode  string
}

type Body struct {
	Data [DATA_SIZE]string
	Desc [DATA_SIZE]string
	// Nonce uint32
	Time string
	Hash string
	Sign string
}
