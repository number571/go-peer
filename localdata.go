package gopeer

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

type conndata struct {
	Certificate string
	Public      string
	Session     string
}

// Get connection and check package.
func runServer(handle func(*Client, *Package), listener *Listener) {
	for {
		if listener.listen == nil {
			break
		}
		conn, err := listener.listen.Accept()
		if err != nil {
			break
		}
		go server(handle, listener, conn)
	}
}

// Package events. By default:
// 1) Check last hash;
// 2) Connection;
// 3) Disconnection;
func server(handle func(*Client, *Package), listener *Listener, conn net.Conn) {
	defer conn.Close()
	for {
		pack := readPackage(conn)
		if pack == nil {
			continue
		}
		pack.receive(handle, listener, conn)
	}
}

func (client *Client) rememberHash(hash string) bool {
	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	if _, ok := client.remember.mapping[hash]; ok {
		return true
	}
	client.remember.index = client.remember.index % settings.REMEMBER
	client.remember.mapping[hash] = client.remember.index
	if _, ok := client.remember.mapping[client.remember.listing[client.remember.index]]; ok {
		delete(client.remember.mapping, client.remember.listing[client.remember.index])
	}
	client.remember.listing[client.remember.index] = hash
	return false
}

// Receive package.
func (pack *Package) receive(handle func(*Client, *Package), listener *Listener, conn net.Conn) {
	if pack.Body.Desc.Redirection >= settings.REDIRECT_QUAN {
		return
	}
	client, ok := listener.Clients[pack.To.Hashname]
	if !ok {
		if pack.To.Hashname == pack.To.Receiver.Hashname {
			return
		}
		client, ok = listener.Clients[pack.To.Receiver.Hashname]
		if !ok {
			return
		}
	}
	if client.rememberHash(pack.Body.Desc.Hash) {
		return
	}
	if pack.To.Hashname != pack.To.Receiver.Hashname {
		if client.InConnections(pack.To.Receiver.Hashname) {
			hash := pack.To.Receiver.Hashname
			pack.To.Hashname = hash
			pack.To.Address = client.Connections[hash].Address
			client.send(RAW, pack)
		} else {
			pack.Body.Desc.Redirection++
			for hash := range client.Connections {
				if hash == pack.From.Sender.Hashname {
					continue
				}
				pack.To.Hashname = hash
				pack.To.Address = client.Connections[hash].Address
				client.send(RAW, pack)
			}
		}
		return
	}

	pack, wasEncrypted := client.tryDecrypt(pack)
	if err := client.isValid(pack); err != nil {
		fmt.Println(err)
		return
	}

	// printJson(pack)

	handleIsUsed := client.HandleAction(settings.TITLE_CONNECT, pack,
		func(client *Client, pack *Package) (set string) {
			client.connectGet(pack, conn)
			return set
		},
		func(client *Client, pack *Package) {
			hash := pack.From.Sender.Hashname
			if !client.InConnections(hash) {
				return
			}
			client.Connections[hash].transfer.isBlocked <- true
		},
	)

	// Subsequent verification is carried out only if the data has been encrypted.
	if handleIsUsed || !wasEncrypted {
		return
	}

	switch pack.Head.Title {
	case settings.TITLE_DISCONNECT:
		switch pack.Head.Option {
		case settings.OPTION_GET:
			client.disconnectGet(pack)
			return
		}
	}

	handleIsUsed = client.HandleAction(settings.TITLE_FILETRANSFER, pack,
		func(client *Client, pack *Package) (set string) {
			nullpack := string(PackJSON(FileTransfer{
				Head: HeadTransfer{
					IsNull: true,
				},
			}))

			if !client.sharing.perm {
				return nullpack
			}

			var read = new(FileTransfer)
			UnpackJSON([]byte(pack.Body.Data), read)

			if read.Head.IsNull {
				return nullpack
			}

			name := strings.Replace(read.Head.Name, "..", "", -1)
			data := readFile(client.sharing.path+name, read.Head.Id)

			return string(PackJSON(FileTransfer{
				Head: HeadTransfer{
					Id:     read.Head.Id,
					Name:   name,
					IsNull: data == nil,
				},
				Body: BodyTransfer{
					Hash: HashSum(data),
					Data: data,
				},
			}))
		},
		func(client *Client, pack *Package) {
			var read = new(FileTransfer)
			UnpackJSON([]byte(pack.Body.Data), read)

			hash := pack.From.Sender.Hashname
			if read.Head.IsNull {
				client.Connections[hash].transfer.isBlocked <- false
				return
			}

			name := client.Connections[hash].transfer.outputFile
			if read.Head.Id == 0 && fileIsExist(name) {
				client.Connections[hash].transfer.isBlocked <- false
				return
			}

			data := read.Body.Data
			if !bytes.Equal(read.Body.Hash, HashSum(data)) {
				client.Connections[hash].transfer.isBlocked <- false
				return
			}

			writeFile(name, read.Body.Data)

			dest := &Destination{
				Address:     client.Connections[pack.From.Hashname].Address,
				Certificate: client.Connections[pack.From.Hashname].Certificate,
				Public:      client.Connections[pack.From.Hashname].Public,
				Receiver:    client.Connections[hash].Public,
			}
			client.SendTo(dest, &Package{
				Head: Head{
					Title:  settings.TITLE_FILETRANSFER,
					Option: settings.OPTION_GET,
				},
				Body: Body{
					Data: string(PackJSON(FileTransfer{
						Head: HeadTransfer{
							Id:   read.Head.Id + 1,
							Name: client.Connections[hash].transfer.inputFile,
						},
					})),
				},
			})
		},
	)

	if handleIsUsed {
		return
	}

	handle(client, pack)
}

// Check package for compliance:
// 1) pack is not null;
// 2) pack.Info.Network == NETWORK;
// 3) pack.Info.Version == VERSION;
// 4) pack.Body.Desc.Difficulty == DIFFICULTY;
// 5) public key can be parsed;
// 6) hash(public) should be equal sender hashname;
// 7) hash(pack) should be equal package hash;
// 8) signature must be created by sender;
// 9) nonce should be equal POW(hash, DIFFICULTY);
func (client *Client) isValid(pack *Package) error {
	if pack == nil {
		return errors.New("pack is null")
	}

	if pack.Info.Network != settings.NETWORK {
		return errors.New("network does not match")
	}

	if pack.Info.Version != settings.VERSION {
		return errors.New("version does not match")
	}

	if pack.Body.Desc.Difficulty != settings.DIFFICULTY {
		return errors.New("difficulty does not match")
	}

	if pack.From.Sender.Hashname == client.Hashname {
		return errors.New("sender and receiver is one person")
	}

	var public *rsa.PublicKey
	if client.InConnections(pack.From.Sender.Hashname) {
		public = client.Connections[pack.From.Sender.Hashname].Public
	} else {
		var data conndata
		UnpackJSON([]byte(pack.Body.Data), &data)
		public = ParsePublic(string(Base64Decode(data.Public)))
	}

	if public == nil {
		return errors.New("can't read public key")
	}

	if HashPublic(public) != pack.From.Sender.Hashname {
		return errors.New("hashname not equal public key")
	}

	hash := HashSum(bytes.Join(
		[][]byte{
			[]byte(pack.Info.Network),
			[]byte(pack.Info.Version),
			[]byte(pack.From.Sender.Hashname),
			[]byte(pack.To.Receiver.Hashname),
			[]byte(pack.Head.Title),
			[]byte(pack.Head.Option),
			[]byte(pack.Body.Data),
			[]byte(pack.Body.Desc.Rand),
		},
		[]byte{},
	))
	if Base64Encode(hash) != pack.Body.Desc.Hash {
		return errors.New("pack hash invalid")
	}

	if Verify(public, hash, Base64Decode(pack.Body.Desc.Sign)) != nil {
		return errors.New("pack sign invalid")
	}

	if !NonceIsValid(Base64Decode(pack.Body.Desc.Hash), uint(pack.Body.Desc.Difficulty), pack.Body.Desc.Nonce) {
		return errors.New("pack nonce is invalid")
	}

	return nil
}

// Send package.
// Check if pack is not null and receive user in saved data.
// Append headers and confirm package.
// Send package.
// If option package is GET, then get response.
// If no response received, then use retrySend() function.
func (client *Client) send(option Option, pack *Package) (*Package, error) {
	switch {
	case pack == nil:
		return nil, errors.New("pack is null")
	case pack.To.Hashname == client.Hashname:
		return nil, errors.New("sender and receiver is one person")
	case !client.InConnections(pack.To.Hashname):
		return nil, errors.New("receiver not in connections")
	}

	pack = client.appendHeaders(pack)
	if option == CONFIRM {
		pack = client.confirmPackage(GenerateRandomBytes(16), pack)
	}
	
	var (
		savedPack = pack
		hash      = pack.To.Hashname
	)

	if client.Connections[hash].Relation == nil {
		ok := client.CertPool.AppendCertsFromPEM([]byte(client.Connections[hash].Certificate))
		if !ok {
			return nil, errors.New("failed to parse root certificate")
		}
		config := &tls.Config{
			ServerName: settings.SERVER_NAME,
			RootCAs:    client.CertPool,
		}
		conn, err := tls.Dial("tcp", pack.To.Address, config)
		if err != nil {
			delete(client.Connections, hash)
			return nil, err
		}
		client.Connections[hash].Relation = conn
		go server(client.listener.handleFunc, client.listener, conn)
	}

	if option == CONFIRM {
		if encPack := client.encryptPackage(pack); encPack != nil {
			pack = encPack
		}
	}

	conn := client.Connections[hash].Relation
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
	if dest.Public == nil && dest.Receiver == nil {
		return nil
	}
	if dest.Receiver == nil {
		dest.Receiver = dest.Public
	}
	hash := HashPublic(dest.Receiver)
	if dest.Public == nil && client.InConnections(hash) {
		dest.Certificate = client.Connections[hash].Certificate
		dest.Public      = client.Connections[hash].ThrowClient
		dest.Address     = client.Connections[hash].Address
	}
	return dest
}

// Connect by GET option.
func (client *Client) connectGet(pack *Package, conn net.Conn) {
	var data conndata
	UnpackJSON([]byte(pack.Body.Data), &data)
	public := ParsePublic(string(Base64Decode(data.Public)))

	hash := pack.From.Sender.Hashname
	client.Connections[hash] = &Connect{
		connected: true,
		transfer: transfer{
			isBlocked: make(chan bool),
		},
		Address:     pack.From.Address,
		Public:      public,
		Certificate: Base64Decode(data.Certificate),
		IsAction:    make(chan bool),
		Session:     DecryptRSA(client.Keys.Private, Base64Decode(data.Session)),
	}

	if pack.From.Hashname == pack.From.Sender.Hashname {
		client.Connections[hash].ThrowClient = public
		client.Connections[hash].Relation = conn
	} else {
		client.Connections[hash].ThrowClient = client.Connections[pack.From.Hashname].Public
	}
}

// Disconnect by GET option.
func (client *Client) disconnectGet(pack *Package) {
	hash := pack.From.Sender.Hashname
	if client.Connections[hash].Relation != nil {
		client.Connections[hash].Relation.Close()
	}
	delete(client.Connections, hash)
}

// Read raw data and convert to package.
func readPackage(conn net.Conn) *Package {
	var (
		message string
		pack    = new(Package)
		size    = uint32(0)
		buffer  = make([]byte, settings.BUFF_SIZE)
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			break
		}
		size += uint32(length)
		if size >= settings.PACK_SIZE {
			return nil
		}
		message += string(buffer[:length])
		// if strings.HasSuffix(message, settings.END_BYTES) {
		if strings.Contains(message, settings.END_BYTES) {
			// message = strings.TrimSuffix(message, settings.END_BYTES)
			message = strings.Split(message, settings.END_BYTES)[0]
			break
		}
	}
	// fmt.Println(size)
	err := json.Unmarshal([]byte(message), pack)
	if err != nil {
		return nil
	}
	return pack
}

// If package not decrypted, then uses first version package.
func (client *Client) tryDecrypt(pack *Package) (*Package, bool) {
	if pack == nil {
		return nil, false
	}
	result := client.decryptPackage(pack)
	if result == nil {
		return pack, false
	}
	return result, true
}

// Encrypt package by session key. Encrypted data:
// 1) Head.Title;
// 2) Head.Option;
// 3) Body.Data;
// 4) Body.Desc;
func (client *Client) encryptPackage(pack *Package) *Package {
	var session []byte

	switch {
	case client.isConnected(pack.To.Receiver.Hashname):
		session = client.Connections[pack.To.Receiver.Hashname].Session
	default:
		return nil
	}

	return &Package{
		Info: Info{
			Network: settings.NETWORK,
			Version: settings.VERSION,
		},
		From: From{
			Sender: Sender{
				Hashname: pack.From.Sender.Hashname,
			},
			Hashname: pack.From.Hashname,
			Address:  pack.From.Address,
		},
		To: To{
			Receiver: Receiver{
				Hashname: pack.To.Receiver.Hashname,
			},
			Hashname: pack.To.Hashname,
			Address:  pack.To.Address,
		},
		Head: Head{
			// Title:  pack.Head.Title,
			// Option: pack.Head.Option,
			Title:  Base64Encode(EncryptAES(session, []byte(pack.Head.Title))),
			Option: Base64Encode(EncryptAES(session, []byte(pack.Head.Option))),
		},
		Body: Body{
			Data: Base64Encode(EncryptAES(session, []byte(pack.Body.Data))),
			Desc: Desc{
				Rand:        Base64Encode(EncryptAES(session, []byte(pack.Body.Desc.Rand))),
				Hash:        pack.Body.Desc.Hash,
				Sign:        pack.Body.Desc.Sign,
				Nonce:       pack.Body.Desc.Nonce,
				Difficulty:  settings.DIFFICULTY,
				Redirection: pack.Body.Desc.Redirection,
			},
		},
	}
}

// Decrypt package by session key. Decrypted data:
// 1) Head.Title;
// 2) Head.Option;
// 3) Body.Data;
// 4) Body.Data.Desc;
func (client *Client) decryptPackage(pack *Package) *Package {
	if !client.InConnections(pack.From.Sender.Hashname) {
		return nil
	}
	session := client.Connections[pack.From.Sender.Hashname].Session
	if DecryptAES(session, Base64Decode(pack.Body.Desc.Rand)) == nil {
		return nil
	}
	return &Package{
		Info: Info{
			Network: settings.NETWORK,
			Version: settings.VERSION,
		},
		From: From{
			Sender: Sender{
				Hashname: pack.From.Sender.Hashname,
			},
			Hashname: pack.From.Hashname,
			Address:  pack.From.Address,
		},
		To: To{
			Receiver: Receiver{
				Hashname: pack.To.Receiver.Hashname,
			},
			Hashname: pack.To.Hashname,
			Address:  pack.To.Address,
		},
		Head: Head{
			// Title:  pack.Head.Title,
			// Option: pack.Head.Option,
			Title:  string(DecryptAES(session, Base64Decode(pack.Head.Title))),
			Option: string(DecryptAES(session, Base64Decode(pack.Head.Option))),
		},
		Body: Body{
			Data: string(DecryptAES(session, Base64Decode(pack.Body.Data))),
			Desc: Desc{
				Rand:        string(DecryptAES(session, Base64Decode(pack.Body.Desc.Rand))),
				Hash:        pack.Body.Desc.Hash,
				Sign:        pack.Body.Desc.Sign,
				Nonce:       pack.Body.Desc.Nonce,
				Difficulty:  settings.DIFFICULTY,
				Redirection: pack.Body.Desc.Redirection,
			},
		},
	}
}

// Get previous hash, generate random bytes, calculate new hash, sign and nonce for package.
// Save current hash in local storage.
func (client *Client) confirmPackage(random []byte, pack *Package) *Package {
	pack.Body.Desc.Rand = Base64Encode(random)
	pack.Body.Desc.Difficulty = settings.DIFFICULTY
	hash := HashSum(bytes.Join(
		[][]byte{
			[]byte(pack.Info.Network),
			[]byte(pack.Info.Version),
			[]byte(pack.From.Sender.Hashname),
			[]byte(pack.To.Receiver.Hashname),
			[]byte(pack.Head.Title),
			[]byte(pack.Head.Option),
			[]byte(pack.Body.Data),
			[]byte(pack.Body.Desc.Rand),
		},
		[]byte{},
	))
	pack.Body.Desc.Hash = Base64Encode(hash)
	pack.Body.Desc.Sign = Base64Encode(Sign(client.Keys.Private, hash))
	pack.Body.Desc.Nonce = ProofOfWork(hash, uint(pack.Body.Desc.Difficulty))
	return pack
}

// Append information about network name, version.
// Append sender information: hashname, public, address.
func (client *Client) appendHeaders(pack *Package) *Package {
	pack.Info.Network = settings.NETWORK
	pack.Info.Version = settings.VERSION
	pack.From.Hashname = client.Hashname
	pack.From.Address = client.Address
	if pack.From.Sender.Hashname == "" {
		pack.From.Sender.Hashname = client.Hashname
	}
	return pack
}

// Check if title is connect.
func (pack *Package) isConnect() bool {
	if pack.Head.Title == settings.TITLE_CONNECT {
		return true
	}
	return false
}

// Check if user connected to client.
func (client *Client) isConnected(hash string) bool {
	if _, ok := client.Connections[hash]; ok {
		return client.Connections[hash].connected
	}
	return false
}

func writeFile(filename string, data []byte) error {
	if !fileIsExist(filename) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(data)
	return nil
}

func readFile(filename string, id uint32) []byte {
	const BEGGINING = 0

	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	_, err = file.Seek(int64(id*settings.FILE_SIZE), BEGGINING)
	if err != nil {
		return nil
	}

	var buffer = make([]byte, settings.FILE_SIZE)
	length, err := file.Read(buffer)
	if err != nil {
		return nil
	}

	return buffer[:length]
}

func fileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// For debug.
func printJson(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(jsonData))
}
