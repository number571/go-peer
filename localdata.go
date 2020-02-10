package gopeer

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type conndata struct {
	Public  string
	Session string
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

	var (
		pack         *Package
		keepconn     bool
		wasEncrypted bool
	)

repeat:
	pack, keepconn = readPackage(conn)
	if pack == nil {
		if keepconn {
			goto repeat
		}
		return
	}

	client, ok := listener.Clients[pack.To.Receiver.Hashname]
	if !ok {
		if pack.To.Hashname == pack.To.Receiver.Hashname {
			if keepconn {
				goto repeat
			}
			return
		}

		client, ok = listener.Clients[pack.To.Hashname]
		if !ok {
			if keepconn {
				goto repeat
			}
			return
		}

		if client.Address == "" {
			if keepconn {
				goto repeat
			}
			return
		}

		client.sendRaw(pack)
		if keepconn {
			goto repeat
		}
		return
	}

	pack, wasEncrypted = client.tryDecrypt(pack)
	// printJson(pack)
	if err := client.isValid(pack); err != nil {
		// fmt.Println(err)
		if keepconn {
			goto repeat
		}
		return
	}

	handleIsUsed := client.HandleAction(settings.TITLE_LASTHASH, pack,
		func(client *Client, pack *Package) (set string) {
			hash := pack.From.Sender.Hashname
			if !client.InConnections(hash) {
				return
			}
			return client.Connections[hash].lastHash
		},
		func(client *Client, pack *Package) {
			hash := pack.From.Sender.Hashname
			if !client.InConnections(hash) {
				return
			}
			client.Connections[hash].lastHash = pack.Body.Data
		},
	)

	if handleIsUsed {
		if keepconn {
			goto repeat
		}
		return
	}

	handleIsUsed = client.HandleAction(settings.TITLE_CONNECT, pack,
		func(client *Client, pack *Package) (set string) {
			client.connectGet(pack, conn)
			return set
		},
		func(client *Client, pack *Package) {
			hash := pack.From.Sender.Hashname
			if !client.InConnections(hash) {
				return
			}
			client.Connections[hash].connected = true
			client.Connections[hash].lastHash = pack.Body.Desc.CurrHash
		},
	)

	// Subsequent verification is carried out only if the data has been encrypted.
	if handleIsUsed || !wasEncrypted {
		if keepconn {
			goto repeat
		}
		return
	}

	hash := pack.From.Sender.Hashname
	client.Connections[hash].lastHash = pack.Body.Desc.CurrHash

	switch pack.Head.Title {
	case settings.TITLE_DISCONNECT:
		switch pack.Head.Option {
		case settings.OPTION_GET:
			client.disconnectGet(pack)
			return
		case settings.OPTION_SET:
			client.Connections[hash].waiting <- true
			delete(client.Connections, hash)
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

			if !client.Sharing.Perm {
				return nullpack
			}

			var read = new(FileTransfer)
			UnpackJSON([]byte(pack.Body.Data), read)

			hash := pack.From.Sender.Hashname
			client.Connections[hash].transfer.isBlocked = true

			if read.Head.IsNull {
				client.Connections[hash].transfer.isBlocked = false
				return nullpack
			}

			name := strings.Replace(read.Head.Name, "..", "", -1)
			data := readFile(client.Sharing.Path+name, read.Head.Id)

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
				client.Connections[hash].transfer.isBlocked = false
				return
			}

			name := client.Connections[hash].transfer.outputFile
			if read.Head.Id == 0 && fileIsExist(name) {
				client.Connections[hash].transfer.isBlocked = false
				return
			}

			data := read.Body.Data
			if !bytes.Equal(read.Body.Hash, HashSum(data)) {
				client.Connections[hash].transfer.isBlocked = false
				return
			}

			writeFile(name, read.Body.Data)
		},
	)

	if handleIsUsed || client.isBlocked(hash) {
		if keepconn {
			goto repeat
		}
		return
	}

	handle(client, pack)
	if keepconn {
		goto repeat
	}
}

// Send raw package without checks, retry's and encryptions.
func (client *Client) sendRaw(pack *Package) (*Package, error) {
	switch {
	case pack == nil: 
		return nil, errors.New("pack is null")
	case pack.To.Receiver.Hashname == client.Hashname: 
		return nil, errors.New("sender and receiver is one person")
	case !client.InConnections(pack.To.Receiver.Hashname): 
		return nil, errors.New("receiver not in connections")
	case client.isBlocked(pack.To.Receiver.Hashname) && pack.Head.Title != settings.TITLE_FILETRANSFER:
		return nil, errors.New("connections is blocked [receiver 2]")
	}
	var (
		savedPack = pack
		hash      = pack.To.Receiver.Hashname
	)

	if client.Connections[hash].Relation != nil {
		conn := client.Connections[hash].Relation
		_, err := conn.Write(
			bytes.Join(
				[][]byte{
					EncryptAES([]byte(settings.NOISE), PackJSON(pack)),
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
	} else {
		conn, err := net.Dial("tcp", pack.To.Address)
		if err != nil {
			delete(client.Connections, hash)
			return nil, err
		}
		conn.Write(EncryptAES([]byte(settings.NOISE), PackJSON(pack)))
		conn.Close()
	}

	return savedPack, nil
}

// Send package.
// Check if pack is not null and receive user in saved data.
// Append headers and confirm package.
// Send package.
// If option package is GET, then get response.
// If no response received, then use retrySend() function.
func (client *Client) send(pack *Package) (*Package, error) {
	switch {
	case pack == nil: 
		return nil, errors.New("pack is null")
	case pack.To.Receiver.Hashname == client.Hashname: 
		return nil, errors.New("sender and receiver is one person")
	case !client.InConnections(pack.To.Hashname): 
		return nil, errors.New("receiver 1 not in connections")
	case !client.InConnections(pack.To.Receiver.Hashname): 
		return nil, errors.New("receiver 2 not in connections")
	case client.isBlocked(pack.To.Hashname) && pack.Head.Title != settings.TITLE_FILETRANSFER:
		return nil, errors.New("connections is blocked [receiver 1]")
	case client.isBlocked(pack.To.Receiver.Hashname) && pack.Head.Title != settings.TITLE_FILETRANSFER:
		return nil, errors.New("connections is blocked [receiver 2]")
	}

	client.confirmPackage(client.appendHeaders(pack))
	var (
		err error
		savedPack = pack
		hash      = pack.To.Receiver.Hashname
	)

	// printJson(pack)
	if encPack := client.encryptPackage(pack); encPack != nil {
		pack = encPack
	}

	if client.Address == "" && client.Connections[hash].Relation == nil {
		conn, err := net.Dial("tcp", pack.To.Address)
		if err != nil {
			delete(client.Connections, hash)
			return nil, err
		}
		client.Connections[hash].Relation = conn
		go server(client.listener.handleFunc, client.listener, conn)
	}

	if client.Connections[hash].Relation != nil {
		conn := client.Connections[hash].Relation
		_, err := conn.Write(
			bytes.Join(
				[][]byte{
					EncryptAES([]byte(settings.NOISE), PackJSON(pack)),
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
	} else {
		conn, err := net.Dial("tcp", pack.To.Address)
		if err != nil {
			delete(client.Connections, hash)
			return nil, err
		}
		conn.Write(EncryptAES([]byte(settings.NOISE), PackJSON(pack)))
		conn.Close()
	}

	if savedPack.Head.Option == settings.OPTION_GET {
		select {
		case <-client.Connections[hash].waiting:
			err = nil
		case <-time.After(time.Duration(settings.RETRY_TIME) * time.Second):
			if savedPack.isConnect() {
				delete(client.Connections, hash)
				return savedPack, fmt.Errorf("time is over (%d seconds)", settings.RETRY_TIME)
			}
			if savedPack.isLasthash() {
				return savedPack, fmt.Errorf("time is over (%d seconds)", settings.RETRY_TIME)
			}
			return savedPack, client.retrySend(savedPack)
		}
	}
	return savedPack, err
}

func (client *Client) isBlocked(hashname string) bool {
	if !client.InConnections(hashname) {
		return false
	}
	if client.Connections[hashname].transfer.isBlocked {
		return true
	}
	return false
}

// Connect by GET option.
func (client *Client) connectGet(pack *Package, conn net.Conn) {
	var (
		public1  *rsa.PublicKey
		data     conndata
		lastHash = ""
	)
	UnpackJSON([]byte(pack.Body.Data), &data)
	public := ParsePublic(string(Base64Decode(data.Public)))

	if pack.From.Hashname == pack.From.Sender.Hashname {
		public1 = public
	} else {
		public1 = client.Connections[pack.From.Hashname].Public
	}

	if client.InConnections(pack.From.Sender.Hashname) {
		lastHash = client.Connections[pack.From.Sender.Hashname].lastHash
	}

	client.Connections[pack.From.Sender.Hashname] = &Connect{
		connected:  true,
		lastHash:   lastHash,
		waiting:    make(chan bool),
		Session:    DecryptRSA(client.Keys.Private, Base64Decode(data.Session)),
		Address:    pack.From.Address,
		Public:     public1,
		PublicRecv: public,
	}

	if pack.From.Address == "" {
		client.Connections[pack.From.Sender.Hashname].Relation = conn
	}

	client.Connections[pack.From.Sender.Hashname].lastHash = pack.Body.Desc.CurrHash
}

// Disconnect by GET option.
func (client *Client) disconnectGet(pack *Package) {
	hash := pack.From.Sender.Hashname
	dest := NewDestination(&Destination{
		Address:  pack.From.Address,
		Public:   client.Connections[hash].Public,
		Receiver: client.Connections[hash].PublicRecv,
	})
	client.SendTo(dest, &Package{
		Head: Head{
			Title:  settings.TITLE_DISCONNECT,
			Option: settings.OPTION_SET,
		},
	})
	if client.Connections[hash].Relation != nil {
		client.Connections[hash].Relation.Close()
	}
	delete(client.Connections, hash)
}

// Read raw data and convert to package.
func readPackage(conn net.Conn) (*Package, bool) {
	var (
		keepconn bool
		message string
		pack    = new(Package)
		size    = uint32(0)
		buffer  = make([]byte, settings.BUFFSIZE)
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			break
		}
		size += uint32(length)
		if size >= settings.PACKSIZE {
			return nil, false
		}
		message += string(buffer[:length])
		// if strings.HasSuffix(message, settings.END_BYTES) {
		if strings.Contains(message, settings.END_BYTES) {
			// message = strings.TrimSuffix(message, settings.END_BYTES)
			message = strings.Split(message, settings.END_BYTES)[0]
			keepconn = true
			break
		}
	}
	// fmt.Println(size)
	err := json.Unmarshal(DecryptAES([]byte(settings.NOISE), []byte(message)), pack)
	if err != nil {
		return nil, false
	}
	return pack, keepconn
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
// 4) Body.Data.Desc;
func (client *Client) encryptPackage(pack *Package) *Package {
	var session []byte

	if pack.isLasthash() {
		return nil
	}

	switch {
	case client.isConnected(pack.To.Receiver.Hashname):
		session = client.Connections[pack.To.Receiver.Hashname].Session
	case client.Connections[pack.To.Hashname].prevSession != nil:
		session = client.Connections[pack.To.Receiver.Hashname].prevSession
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
				Rand:       Base64Encode(EncryptAES(session, []byte(pack.Body.Desc.Rand))),
				PrevHash:   pack.Body.Desc.PrevHash,
				CurrHash:   pack.Body.Desc.CurrHash,
				Sign:       pack.Body.Desc.Sign,
				Nonce:      pack.Body.Desc.Nonce,
				Difficulty: settings.DIFFICULTY,
			},
		},
	}
}

// Check if user connected to client.
func (client *Client) isConnected(hash string) bool {
	if _, ok := client.Connections[hash]; ok {
		return client.Connections[hash].connected
	}
	return false
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
				Rand:       string(DecryptAES(session, Base64Decode(pack.Body.Desc.Rand))),
				PrevHash:   pack.Body.Desc.PrevHash,
				CurrHash:   pack.Body.Desc.CurrHash,
				Sign:       pack.Body.Desc.Sign,
				Nonce:      pack.Body.Desc.Nonce,
				Difficulty: settings.DIFFICULTY,
			},
		},
	}
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
// 10) IF sender in connections and package is not LASTHASH
// and package is not have OPTION_SET:
// 10.1) saved last hash should be equal previous package hash;
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
		public = client.Connections[pack.From.Sender.Hashname].PublicRecv
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

	tempPack := *pack

	tempPack.Body.Desc.CurrHash = ""
	tempPack.Body.Desc.Sign = ""
	tempPack.Body.Desc.Nonce = 0
	tempPack.Body.Desc.Difficulty = settings.DIFFICULTY

	hash := HashSum(PackJSON(&tempPack))
	if Base64Encode(hash) != pack.Body.Desc.CurrHash {
		return errors.New("pack hash invalid")
	}

	if Verify(public, hash, Base64Decode(pack.Body.Desc.Sign)) != nil {
		return errors.New("pack sign invalid")
	}

	if !NonceIsValid(Base64Decode(pack.Body.Desc.CurrHash), uint(pack.Body.Desc.Difficulty), pack.Body.Desc.Nonce) {
		return errors.New("pack nonce is invalid")
	}

	if client.InConnections(pack.From.Sender.Hashname) {
		if pack.isLasthash() {
			return nil
		}

		if client.Connections[pack.From.Sender.Hashname].lastHash != pack.Body.Desc.PrevHash {
			return errors.New("prev pack hash not equal hash in saved")
		}
	}

	return nil
}

// Append information about network name, version.
// Append sender information: hashname, public, address.
func (client *Client) appendHeaders(pack *Package) *Package {
	pack.Info.Network = settings.NETWORK
	pack.Info.Version = settings.VERSION
	pack.From.Hashname = client.Hashname
	pack.From.Address = client.Address

	if pack.From.Sender.Hashname == "" {
		pack.From.Sender.Hashname = pack.From.Hashname
	}

	if pack.To.Hashname != pack.To.Receiver.Hashname {
		pack.From.Address = pack.To.Address
		pack.From.Hashname = pack.To.Hashname
	}
	return pack
}

// Get previous hash, generate random bytes, calculate new hash, sign and nonce for package.
// Save current hash in local storage.
func (client *Client) confirmPackage(pack *Package) *Package {
	if !pack.isLasthash() {
		pack.Body.Desc.PrevHash = client.Connections[pack.To.Receiver.Hashname].lastHash
	}

	pack.Body.Desc.Rand = Base64Encode(GenerateRandomBytes(16))
	pack.Body.Desc.CurrHash = ""
	pack.Body.Desc.Sign = ""
	pack.Body.Desc.Nonce = 0
	pack.Body.Desc.Difficulty = settings.DIFFICULTY

	hash := HashSum(PackJSON(pack))
	pack.Body.Desc.CurrHash = Base64Encode(hash)
	pack.Body.Desc.Sign = Base64Encode(Sign(client.Keys.Private, hash))
	pack.Body.Desc.Nonce = ProofOfWork(hash, uint(pack.Body.Desc.Difficulty))

	if !pack.isLasthash() { // && pack.Head.Option == settings.OPTION_SET
		client.Connections[pack.To.Receiver.Hashname].lastHash = pack.Body.Desc.CurrHash
	}

	return pack
}

// Check if is package to request the last hash.
func (pack *Package) isLasthash() bool {
	if pack.Head.Title == settings.TITLE_LASTHASH {
		return true
	}
	return false
}

// Check if title is connect.
func (pack *Package) isConnect() bool {
	if pack.Head.Title == settings.TITLE_CONNECT {
		return true
	}
	return false
}

// Retry sending the package RETRY_NUMB times.
// Inerval equal RETRY_TIME.
func (client *Client) retrySend(pack *Package) error {
	if pack == nil {
		return errors.New("pack is null")
	}
	if !client.InConnections(pack.To.Receiver.Hashname) {
		return errors.New("receiver not in connections")
	}
	var (
		retryNum = settings.RETRY_NUMB
		hash     = pack.To.Receiver.Hashname
	)

retry:
	client.Connections[hash].lastHash = ""
	dest := NewDestination(&Destination{
		Address: pack.To.Address,
		Public:  client.Connections[hash].Public,
		Receiver:  client.Connections[hash].PublicRecv,
	})
	_, err := client.SendTo(dest, &Package{
		Head: Head{
			Title:  settings.TITLE_LASTHASH,
			Option: settings.OPTION_GET,
		},
	})
	if err != nil {
		if retryNum > 1 {
			retryNum--
			goto retry
		}
		// delete(client.Connections, hash)
		return fmt.Errorf("time is over (%d seconds)", settings.RETRY_TIME)
	}

	client.confirmPackage(client.appendHeaders(pack))
	var (
		savedPack = pack
	)
	if encPack := client.encryptPackage(pack); encPack != nil {
		pack = encPack
	}

	if client.Connections[hash].Relation != nil {
		conn := client.Connections[hash].Relation
		_, err := conn.Write(
			bytes.Join(
				[][]byte{
					EncryptAES([]byte(settings.NOISE), PackJSON(pack)),
					[]byte(settings.END_BYTES),
				},
				[]byte{},
			),
		)
		if err != nil {
			conn.Close()
			delete(client.Connections, hash)
			return err
		}
	} else {
		conn, err := net.Dial("tcp", pack.To.Address)
		if err != nil {
			delete(client.Connections, hash)
			return err
		}
		conn.Write(EncryptAES([]byte(settings.NOISE), PackJSON(pack)))
		conn.Close()
	}

	if savedPack.Head.Option == settings.OPTION_GET {
		select {
		case <-client.Connections[hash].waiting:
			err = nil
		case <-time.After(time.Duration(settings.RETRY_TIME) * time.Second):
			pack = savedPack
			if retryNum > 1 {
				retryNum--
				goto retry
			}
			// delete(client.Connections, hash)
			return fmt.Errorf("time is over (%d seconds)", settings.RETRY_TIME)
		}
	}
	return err
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

	_, err = file.Seek(int64(id*settings.FILESIZE), BEGGINING)
	if err != nil {
		return nil
	}

	var buffer = make([]byte, settings.FILESIZE)
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
