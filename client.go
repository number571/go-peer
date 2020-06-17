package gopeer

import (
	"net"
	"bytes"
	"crypto/rsa"
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

// Return listener certificate.
func (client *Client) Certificate() []byte {
	return client.listener.certificate
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

// Wrap function in mutex Lock/Unlock.
func (client *Client) Action(action func()) {
	client.mutex.Lock()
	action()
	client.mutex.Unlock()
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
	if dest == nil {
		return errors.New("dest is null")
	}
	hash := HashPublic(dest.Receiver)
	if !client.InConnections(hash) {
		return errors.New("client not connected")
	}
	_, err = client.SendTo(dest, &Package{
		Head: Head{
			Title:  settings.TITLE_DISCONNECT,
			Option: settings.OPTION_GET,
		},
	})
	if err != nil {
		client.mutex.Lock()
		client.disconnect(hash)
		client.mutex.Unlock()
		return err
	}
	select {
	case <-client.Connections[hash].action:
		// disconnect in OPTION_SET
	case <-time.After(time.Duration(settings.WAITING_TIME) * time.Second):
		client.mutex.Lock()
		client.disconnect(hash)
		client.mutex.Unlock()
	}
	return nil
}

// Connect to user.
// Create local data with parameters.
// After sending GET and receiving SET package, set Connected = true.
func (client *Client) Connect(dest *Destination) error {
	dest = client.wrapDest(dest)
	if dest == nil {
		return errors.New("dest is null")
	}
	var (
		relation net.Conn
		session = GenerateRandomBytes(uint(settings.SESSION_SIZE))
		hash    = HashPublic(dest.Receiver)
	)
	if client.InConnections(hash) { // dest.Address == settings.IS_CLIENT && 
		relation = client.Connections[hash].relation
	}
	if dest.Public == nil {
		return client.hiddenConnect(hash, session, dest.Receiver, relation)
	}
	client.mutex.Lock()
	client.Connections[hash] = &Connect{
		connected:   false,
		hashname:    hash,
		address:     dest.Address,
		throwClient: dest.Public,
		public:      dest.Receiver,
		certificate: dest.Certificate,
		session:     session,
		relation:    relation,
		action:      make(chan bool),
		Action:      make(chan bool),
	}
	client.mutex.Unlock()
	var count = settings.RETRY_QUAN
repeat:
	_, err := client.SendTo(dest, &Package{
		Head: Head{
			Title:  settings.TITLE_CONNECT,
			Option: settings.OPTION_GET,
		},
		Body: Body{
			Data: string(PackJSON(conndata{
				Certificate: Base64Encode(client.listener.certificate),
				Public:      Base64Encode([]byte(StringPublic(client.keys.public))),
				Session:     Base64Encode(EncryptRSA(dest.Receiver, session)),
			})),
		},
	})
	if err != nil {
		return err
	}
	select {
	case <-client.Connections[hash].action:
		client.mutex.Lock()
		client.Connections[hash].connected = true
		client.mutex.Unlock()
	case <-time.After(time.Duration(settings.WAITING_TIME) * time.Second):
		if count > 0 {
			count--
			goto repeat
		}
		client.mutex.Lock()
		client.disconnect(hash)
		client.mutex.Unlock()
		return errors.New("client not connected")
	}
	return nil
}

// Load file from node.
// Input = name file in node side.
// Output = result name file in our side.
func (client *Client) LoadFile(dest *Destination, input string, output string) error {
	dest = client.wrapDest(dest)
	if dest == nil {
		return errors.New("dest is null")
	}
	hash := HashPublic(dest.Receiver)
	if !client.InConnections(hash) {
		return errors.New("client not connected")
	}
	if fileIsExist(output) {
		return errors.New("file already exists")
	}
	client.mutex.Lock()
	client.Connections[hash].transfer.active = true
	client.mutex.Unlock()
	defer func() {
		client.mutex.Lock()
		client.Connections[hash].transfer.active = false
		client.mutex.Unlock()
	}()
	var (
		read  = new(FileTransfer)
		count = settings.RETRY_QUAN
	)
	for id := uint64(0); ; id++ {
	repeat:
		_, err := client.SendTo(dest, &Package{
			Head: Head{
				Title:  settings.TITLE_FILETRANSFER,
				Option: settings.OPTION_GET,
			},
			Body: Body{
				Data: string(PackJSON(FileTransfer{
					Head: HeadTransfer{
						Id:   id,
						Name: input,
					},
				})),
			},
		})
		if err != nil {
			return err
		}
		select {
		case <-client.Connections[hash].action:
			// pass
		case <-time.After(time.Duration(settings.WAITING_TIME*2) * time.Second):
			if count > 0 {
				count--
				goto repeat
			}
			return errors.New("waiting time is over")
		}
		UnpackJSON([]byte(client.Connections[hash].transfer.packdata), read)
		if read == nil {
			return errors.New("pack is null")
		}
		if read.Head.IsNull {
			break
		}
		if read.Head.Id != id {
			return errors.New("id not equal file part id")
		}
		if read.Head.Name != input {
			return errors.New("input name not equal file part name")
		}
		if !bytes.Equal(read.Body.Hash, HashSum(read.Body.Data)) {
			return errors.New("hash not equal file part hash")
		}
		writeFile(output, read.Body.Data)
		count = settings.RETRY_QUAN
		
	}
	return nil
}

// Send by Destination.
func (client *Client) SendTo(dest *Destination, pack *Package) (*Package, error) {
	dest = client.wrapDest(dest)
	switch {
	case dest == nil:
		return nil, errors.New("dest is null")
	case dest.Public == nil:
		return nil, errors.New("public is null")
	case dest.Receiver == nil:
		return nil, errors.New("receiver is null")
	}
	pack.To.Receiver.Hashname = HashPublic(dest.Receiver)
	pack.To.Hashname = HashPublic(dest.Public)
	pack.To.Address = dest.Address
	return client.send(_confirm, pack)
}
