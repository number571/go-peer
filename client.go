package gopeer

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"errors"
	"time"
)

// Return client hashname.
func (client *Client) Hashname() string {
	return client.hashname
}

// Return client's public key.
func (client *Client) Public() *rsa.PublicKey {
	x := *client.keys.public
	return &x
}

// Return client's private key.
func (client *Client) Private() *rsa.PrivateKey {
	x := *client.keys.private
	return &x
}

// Return listener address.
func (client *Client) Address() string {
	return client.address
}

// Return Destination struct from connected client.
func (client *Client) Destination(hash string) *Destination {
	if !client.InConnections(hash) {
		return nil
	}
	return &Destination{
		Address:     client.Connections[hash].address,
		Certificate: client.Connections[hash].certificate,
		Public:      client.Connections[hash].throwClient,
		Receiver:    client.Connections[hash].public,
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
		connected:   false,
		hashname:    hash,
		address:     dest.Address,
		throwClient: dest.Public,
		public:      dest.Receiver,
		certificate: dest.Certificate,
		session:     session,
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
				Public:  Base64Encode([]byte(StringPublic(client.keys.public))),
				Session: Base64Encode(EncryptRSA(dest.Receiver, session)),
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

// Find hidden connection throw nodes.
func (client *Client) hiddenConnect(hash string, session []byte, receiver *rsa.PublicKey) error {
	var (
		random = GenerateRandomBytes(16)
		pack   = &Package{
			Head: Head{
				Title:  settings.TITLE_CONNECT,
				Option: settings.OPTION_GET,
			},
			Body: Body{
				Data: string(PackJSON(conndata{
					Certificate: Base64Encode(client.listener.certificate),
					Public:      Base64Encode([]byte(StringPublic(client.keys.public))),
					Session:     Base64Encode(EncryptRSA(receiver, session)),
				})),
			},
		}
	)
	for _, conn := range client.Connections {
		client.Connections[hash] = &Connect{
			connected: false,
			Chans: Chans{
				Action: make(chan bool),
				action: make(chan bool),
			},
			address:     conn.address,
			throwClient: conn.public,
			public:      receiver,
			certificate: conn.certificate,
			session:     session,
		}
		pack.To.Receiver.Hashname = hash
		pack.To.Hashname = HashPublic(conn.public)
		pack.To.Address = conn.address
		pack = client.confirmPackage(random, client.appendHeaders(pack))
		_, err := client.send(_raw, pack)
		if err != nil {
			continue
		}
		select {
		case <-client.Connections[hash].Chans.action:
			client.Connections[hash].connected = true
			return nil
		case <-time.After(time.Duration(settings.WAITING_TIME) * time.Second):
			if client.Connections[hash].relation != nil {
				client.Connections[hash].relation.Close()
			}
			delete(client.Connections, hash)
		}
	}
	return errors.New("Connection undefined")
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
	defer func() {
		client.Connections[hash].transfer.active = false
	}()

	for i := uint32(0) ;; i++ {
		client.SendTo(dest, &Package{
			Head: Head{
				Title:  settings.TITLE_FILETRANSFER,
				Option: settings.OPTION_GET,
			},
			Body: Body{
				Data: string(PackJSON(FileTransfer{
					Head: HeadTransfer{
						Id:   i,
						Name: input,
					},
				})),
			},
		})

		select {
		case <-client.Connections[hash].Chans.action:
			// pass
		case <-time.After(time.Duration(settings.WAITING_TIME) * time.Second):
			return errors.New("waiting time is over")
		}

		var read = new(FileTransfer)
		UnpackJSON([]byte(client.Connections[hash].transfer.packdata), read)

		if read == nil {
			return errors.New("pack is null")
		}

		if read.Head.IsNull {
			break
		}

		if read.Head.Id == 0 && fileIsExist(output) {
			return errors.New("file already exists")
		}

		data := read.Body.Data
		if !bytes.Equal(read.Body.Hash, HashSum(data)) {
			return errors.New("hash not equal file hash")
		}

		writeFile(output, read.Body.Data)
	}

	return nil
}

// Send by Destination.
func (client *Client) SendTo(dest *Destination, pack *Package) (*Package, error) {
	dest = client.wrapDest(dest)
	if dest == nil {
		return nil, errors.New("dest is null")
	}

	pack.To.Receiver.Hashname = HashPublic(dest.Receiver)
	pack.To.Hashname = HashPublic(dest.Public)
	pack.To.Address = dest.Address

	return client.send(_confirm, pack)
}

// Send package.
// Check if pack is not null and receive user in saved data.
// Append headers and confirm package.
// Send package.
// If option package is GET, then get response.
// If no response received, then use retrySend() function.
func (client *Client) send(option optionType, pack *Package) (*Package, error) {
	switch {
	case pack == nil:
		return nil, errors.New("pack is null")
	case pack.To.Hashname == client.hashname:
		return nil, errors.New("sender and receiver is one person")
	case !client.InConnections(pack.To.Hashname):
		return nil, errors.New("receiver not in connections")
	}

	pack = client.appendHeaders(pack)
	if option == _confirm {
		pack = client.confirmPackage(GenerateRandomBytes(16), pack)
	}

	var (
		savedPack = pack
		hash      = pack.To.Hashname
	)

	if client.Connections[hash].relation == nil {
		ok := client.certPool.AppendCertsFromPEM([]byte(client.Connections[hash].certificate))
		if !ok {
			return nil, errors.New("failed to parse root certificate")
		}
		config := &tls.Config{
			ServerName: settings.SERVER_NAME,
			RootCAs:    client.certPool,
		}
		conn, err := tls.Dial("tcp", pack.To.Address, config)
		if err != nil {
			delete(client.Connections, hash)
			return nil, err
		}
		client.Connections[hash].relation = conn
		go serveClient(client.listener, client, client.listener.handleFunc, hash, conn)
	}

	if option == _confirm {
		if encPack := client.encryptPackage(pack); encPack != nil {
			pack = encPack
		}
	}

	conn := client.Connections[hash].relation
	_, err := conn.Write(
		bytes.Join(
			[][]byte{
				PackJSON(pack),
				[]byte(settings.END_BYTES),
			},
			[]byte{},
		),
	)
	if err != nil {
		conn.Close()
		delete(client.Connections, hash)
		return nil, err
	}

	return savedPack, err
}

func (client *Client) wrapDest(dest *Destination) *Destination {
	if dest == nil {
		return nil
	}
	if dest.Public == nil && dest.Receiver == nil {
		return nil
	}
	if dest.Receiver == nil {
		dest.Receiver = dest.Public
	}
	hash := HashPublic(dest.Receiver)
	if dest.Public == nil && client.InConnections(hash) {
		dest.Certificate = client.Connections[hash].certificate
		dest.Public = client.Connections[hash].throwClient
		dest.Address = client.Connections[hash].address
	}
	return dest
}
