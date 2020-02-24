package gopeer

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"strings"
	"sync"
	"time"
)

// Create new Listener by address "ipv4:port".
func NewListener(address string) *Listener {
	if address == settings.IS_CLIENT {
		return &Listener{
			Address: Address{
				Ipv4: settings.IS_CLIENT,
			},
			Clients: make(map[string]*Client),
		}
	}
	splited := strings.Split(address, ":")
	if len(splited) != 2 {
		return nil
	}
	return &Listener{
		Address: Address{
			Ipv4: splited[0],
			Port: ":" + splited[1],
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
		Hashname: hash,
		Keys: Keys{
			Private: private,
			Public:  public,
		},
		Mutex:       new(sync.Mutex),
		Address:     listener.Address.Ipv4 + listener.Address.Port,
		CertPool:    x509.NewCertPool(),
		Connections: make(map[string]*Connect),
	}
	return listener.Clients[hash]
}

// Open connection for receiving data.
func (listener *Listener) Open(c *Certificate) *Listener {
	cert, err := tls.X509KeyPair(c.Cert, c.Key)
	if err != nil {
		return nil
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	if listener.Address.Ipv4+listener.Address.Port == settings.IS_CLIENT {
		return listener
	}
	listener.Certificate = c.Cert
	listener.listen, err = tls.Listen("tcp", settings.TEMPLATE+listener.Address.Port, config)
	if err != nil {
		return nil
	}
	return listener
}

// Run handle server for listening packages.
func (listener *Listener) Run(handle func(*Client, *Package)) *Listener {
	listener.handleFunc = handle
	if listener.Address.Ipv4+listener.Address.Port == settings.IS_CLIENT {
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
	if listener.listen == nil {
		return
	}
	listener.listen.Close()
}

func (client *Client) Destination(hash string) *Destination {
	if !client.InConnections(hash) {
		return nil
	}
	return &Destination{
		Address:     client.Connections[hash].Address,
		Certificate: client.Connections[hash].Certificate,
		Public:      client.Connections[hash].ThrowClient,
		Receiver:    client.Connections[hash].Public,
	}
}

// Check if user saved in client data.
func (client *Client) InConnections(hash string) bool {
	if _, ok := client.Connections[hash]; ok {
		return true
	}
	return false
}

// Switcher function about GET and SET options.
// GET: accept package and send response;
// SET: accept package;
func (client *Client) HandleAction(title string, pack *Package, handleGet func(*Client, *Package) string, handleSet func(*Client, *Package)) bool {
	if pack.Head.Title != title {
		return false
	}
	switch pack.Head.Option {
	case settings.OPTION_GET:
		data := handleGet(client, pack)
		hash := pack.From.Sender.Hashname
		client.SendTo(client.Destination(hash), &Package{
			Head: Head{
				Title:  title,
				Option: settings.OPTION_SET,
			},
			Body: Body{
				Data: data,
			},
		})
	case settings.OPTION_SET:
		handleSet(client, pack)
	default: 
		return false
	}
	return true
}

// Disconnect from user.
// Send package for disconnect.
// If the user is not responding: delete in local data.
func (client *Client) Disconnect(dest *Destination) error {
	var err error
	dest = client.wrapDest(dest)

	hash := HashPublic(dest.Receiver)
	if !client.InConnections(hash) {
		return errors.New("client not connected")
	}

	if client.Connections[hash].relation == nil {
		_, err = client.SendTo(dest, &Package{
			Head: Head{
				Title:  settings.TITLE_DISCONNECT,
				Option: settings.OPTION_GET,
			},
		})
	}

	if client.Connections[hash].relation != nil {
		client.Connections[hash].relation.Close()
	}

	delete(client.Connections, hash)
	return err
}

// Connect to user.
// Create local data with parameters.
// After sending GET and receiving SET package, set Connected = true.
func (client *Client) Connect(dest *Destination) error {
	dest = client.wrapDest(dest)
	var (
		session = GenerateRandomBytes(32)
		hash    = HashPublic(dest.Receiver)
	)
	if dest.Public == nil {
		return client.hiddenConnect(hash, session, dest.Receiver)
	}
	client.Connections[hash] = &Connect{
		connected: false,
		Address:     dest.Address,
		ThrowClient: dest.Public,
		Public:      dest.Receiver,
		Certificate: dest.Certificate,
		Session:     session,
		Chans: Chans{
			Action: make(chan bool),
			action: make(chan bool),
		},
	}
	_, err := client.SendTo(dest, &Package{
		Head: Head{
			Title:  settings.TITLE_CONNECT,
			Option: settings.OPTION_GET,
		},
		Body: Body{
			Data: string(PackJSON(conndata{
				Public:      Base64Encode([]byte(StringPublic(client.Keys.Public))),
				Session:     Base64Encode(EncryptRSA(dest.Receiver, session)),
			})),
		},
	})
	if err != nil {
		return err
	}
	select {
	case <-client.Connections[hash].Chans.action:
		client.Connections[hash].connected = true
	case <-time.After(time.Duration(settings.WAITING_TIME) * time.Second):
		if client.Connections[hash].relation != nil {
			client.Connections[hash].relation.Close()
		}
		delete(client.Connections, hash)
		return errors.New("client not connected")
	}
	return nil
}

// Load file from node.
// Input = name file in node side.
// Output = result name file in our side.
func (client *Client) LoadFile(dest *Destination, input string, output string) error {
	dest = client.wrapDest(dest)

	hash := HashPublic(dest.Receiver)
	if !client.InConnections(hash) {
		return errors.New("client not connected")
	}

	client.Connections[hash].transfer.active = true
	client.Connections[hash].transfer.inputFile = input
	client.Connections[hash].transfer.outputFile = output

	client.SendTo(dest, &Package{
		Head: Head{
			Title:  settings.TITLE_FILETRANSFER,
			Option: settings.OPTION_GET,
		},
		Body: Body{
			Data: string(PackJSON(FileTransfer{
				Head: HeadTransfer{
					Id:   0,
					Name: input,
				},
			})),
		},
	})

	<-client.Connections[hash].Chans.action
	client.Connections[hash].transfer.active = false
	return nil
}

// Send by Destination.
func (client *Client) SendTo(dest *Destination, pack *Package) (*Package, error) {
	dest = client.wrapDest(dest)

	pack.To.Receiver.Hashname = HashPublic(dest.Receiver)
	pack.To.Hashname = HashPublic(dest.Public)
	pack.To.Address = dest.Address

	return client.send(CONFIRM, pack)
}
