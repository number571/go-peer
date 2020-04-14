package gopeer

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"strings"
	"sync"
)

// Create new Listener by address "ipv4:port".
func NewListener(addr string) *Listener {
	if addr == settings.IS_CLIENT {
		return &Listener{
			address: address{
				ipv4: settings.IS_CLIENT,
			},
			Clients: make(map[string]*Client),
		}
	}
	splited := strings.Split(addr, ":")
	if len(splited) != 2 {
		return nil
	}
	return &Listener{
		address: address{
			ipv4: splited[0],
			port: ":" + splited[1],
		},
		Clients: make(map[string]*Client),
	}
}

// Create new client in listener by private key.
func (listener *Listener) NewClient(private *rsa.PrivateKey) *Client {
	public := &private.PublicKey
	hash := HashPublic(public)
	listener.Clients[hash] = &Client{
		listener: listener,
		remember: remember{
			mapping: make(map[string]uint16),
			listing: make([]string, settings.REMEMBER),
		},
		F2F: F2F{
			Friends: make(map[string]bool),
		},
		hashname: hash,
		keys: keys{
			private: private,
			public:  public,
		},
		Mutex:       new(sync.Mutex),
		address:     listener.address.ipv4 + listener.address.port,
		certPool:    x509.NewCertPool(),
		Connections: make(map[string]*Connect),
	}
	return listener.Clients[hash]
}

// Open connection for receiving data.
func (listener *Listener) Open(c *Certificate) *Listener {
	if c == nil {
		return nil
	}
	cert, err := tls.X509KeyPair(c.Cert, c.Key)
	if err != nil {
		return nil
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener.certificate = c.Cert
	if listener.address.ipv4+listener.address.port == settings.IS_CLIENT {
		return listener
	}
	listener.listen, err = tls.Listen("tcp", settings.TEMPLATE+listener.address.port, config)
	if err != nil {
		return nil
	}
	return listener
}

// Run handle server for listening packages.
func (listener *Listener) Run(handle func(*Client, *Package)) *Listener {
	listener.handleFunc = handle
	if listener.address.ipv4+listener.address.port == settings.IS_CLIENT {
		return listener
	}
	go runServer(listener, handle)
	return listener
}

// Close listener connection.
func (listener *Listener) Close() {
	if listener == nil {
		return
	}
	
	for i := range listener.Clients {
		for hash := range listener.Clients[i].Connections {
			if listener.Clients[i].Connections[hash].relation == nil {
				continue
			}
			listener.Clients[i].Connections[hash].relation.Close()
		}
	}
	
	if listener.listen != nil {
		listener.listen.Close()
	}
}

// Return listener certificate.
func (listener *Listener) Certificate() []byte {
	return listener.certificate
}

// Return listener address.
func (listener *Listener) Address() string {
	return listener.address.ipv4 + listener.address.port
}
