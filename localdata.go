package gopeer

import (
	"bytes"
	"time"
	"math/big"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

// Receive package.
func (pack *Package) receive(listener *Listener, handle func(*Client, *Package), conn net.Conn) bool {
	if pack.Body.Desc.Redirection >= settings.REDIRECT_QUAN {
		return false
	}
	client, ok := listener.Clients[pack.To.Hashname]
	if !ok {
		return false
	}
	if pack.To.Hashname != pack.To.Receiver.Hashname {
		if err := client.IsValidRedirect(pack); err != nil {
			return false
		}
		if !client.rememberHash(pack.Body.Desc.Hash) {
			return false
		}
		if client.InConnections(pack.To.Receiver.Hashname) {
			hash := pack.To.Receiver.Hashname
			pack.To.Hashname = hash
			pack.To.Address = client.Connections[hash].address
			client.send(_raw, pack)
		} else {
			pack.Body.Desc.Redirection++
			for hash, cli := range client.Connections {
				if hash == pack.From.Sender.Hashname || hash == pack.From.Hashname {
					continue
				}
				pack.To.Hashname = hash
				pack.To.Address = cli.address
				client.send(_raw, pack)
			}
		}
		return false
	}

	pack, wasEncrypted := client.tryDecrypt(pack)
	if err := client.isValid(pack); err != nil {
		// fmt.Println(err)
		return false
	}

	// printJSON(pack)

	if !client.rememberHash(pack.Body.Desc.Hash) {
		return false
	}

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
			client.Connections[hash].Chans.action <- true
		},
	)

	if handleIsUsed {
		return true
	}

	// Subsequent verification is carried out only if the data has been encrypted.
	if !wasEncrypted {
		return false
	}

	handleIsUsed = client.HandleAction(settings.TITLE_DISCONNECT, pack,
		func(client *Client, pack *Package) (set string) {
			client.disconnectGet(pack)
			return set
		},
		func(client *Client, pack *Package) {
		},
	)

	if handleIsUsed {
		return true
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

			if read.Head.IsNull {
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
			hash := pack.From.Sender.Hashname
			if !client.Connections[hash].transfer.active {
				return
			}
			client.Connections[hash].Chans.action <- true
			client.Connections[hash].transfer.packdata = pack.Body.Data
		},
	)

	if handleIsUsed {
		return true
	}

	handle(client, pack)
	return true
}

func (client *Client) IsValidRedirect(pack *Package) error {
	if !client.InConnections(pack.From.Hashname) {
		return errors.New("not in connections")
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
			ToBytes(pack.Body.Desc.Id),
			[]byte(pack.Body.Desc.Rand),
		},
		[]byte{},
	))

	if Base64Encode(hash) != pack.Body.Test.Hash {
		return errors.New("pack hash invalid")
	}

	public := client.Connections[pack.From.Hashname].public
	if Verify(public, hash, Base64Decode(pack.Body.Test.Sign)) != nil {
		return errors.New("pack sign invalid")
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
			hashname:    hash,
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
			delete(client.Connections, hash)
			return nil, errors.New("failed to parse root certificate")
		}
		config := &tls.Config{
			ServerName: settings.NETWORK,
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
	pack = client.signPackage(pack)

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

// Remember package hash.
func (client *Client) rememberHash(hash string) bool {
	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	if _, ok := client.remember.mapping[hash]; ok {
		return false
	}
	client.remember.index = client.remember.index % settings.REMEMBER
	client.remember.mapping[hash] = client.remember.index
	if _, ok := client.remember.mapping[client.remember.listing[client.remember.index]]; ok {
		delete(client.remember.mapping, client.remember.listing[client.remember.index])
	}
	client.remember.listing[client.remember.index] = hash
	return true
}

// Check package for compliance:
// 1) pack is not null;
// 2) pack.Info.Network == NETWORK;
// 3) pack.Info.Version == VERSION;
// 4) pack.Body.Desc.Difficulty == DIFFICULTY;
// 5) public key can be parsed;
// 6) hash(public) should be equal sender hashname;
// 7) check key size;
// 8) check certificate size;
// 9) hash(pack) should be equal package hash;
// 10) signature must be created by sender;
// 11) nonce should be equal POW(hash, DIFFICULTY);
// 12) check package id;
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

	if pack.From.Sender.Hashname == client.hashname {
		return errors.New("sender and receiver is one person")
	}

	if _, ok := client.F2F.Friends[pack.From.Sender.Hashname]; client.F2F.Perm && !ok {
		return errors.New("hashname undefined in list of friends")
	}

	var (
		public *rsa.PublicKey
		certif *x509.Certificate
	)
	if client.InConnections(pack.From.Sender.Hashname) {
		public = client.Connections[pack.From.Sender.Hashname].public
		certif = ParseCertificate(string(client.Connections[pack.From.Sender.Hashname].certificate))
	} else {
		var data conndata
		UnpackJSON([]byte(pack.Body.Data), &data)
		public = ParsePublic(string(Base64Decode(data.Public)))
		certif = ParseCertificate(string(Base64Decode(data.Certificate)))
	}

	if public == nil {
		return errors.New("can't read public key")
	}

	if certif == nil {
		return errors.New("can't read certificate")
	}

	if HashPublic(public) != pack.From.Sender.Hashname {
		return errors.New("hashname not equal public key")
	}

	x := big.NewInt(1)
    x.Lsh(x, uint(settings.KEY_SIZE - 1))
    if public.N.Cmp(x) == -1 {
    	return errors.New("public size < setting key size")
    }

    x = big.NewInt(1)
    x.Lsh(x, uint(settings.KEY_SIZE - 1))
    if certif.PublicKey.(*rsa.PublicKey).N.Cmp(x) == -1 {
    	return errors.New("certificate size < setting cert size")
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
			ToBytes(pack.Body.Desc.Id),
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

	if client.InConnections(pack.From.Sender.Hashname) {
		if pack.Head.Title == settings.TITLE_FILETRANSFER {
			goto pass
		}
		if client.Connections[pack.From.Sender.Hashname].packageId >= settings.max_id && pack.Head.Option == settings.OPTION_SET {
			client.Connect(client.Destination(pack.From.Sender.Hashname))
			return nil
		}
	pass:
		if pack.Head.Title != settings.TITLE_CONNECT && pack.Body.Desc.Id+1 < client.Connections[pack.From.Sender.Hashname].packageId {
			return errors.New("package id < saved package id")
		}
		client.Connections[pack.From.Sender.Hashname].packageId = pack.Body.Desc.Id+1
	}

	return nil
}

// Connect by GET option.
func (client *Client) connectGet(pack *Package, conn net.Conn) {
	var data conndata
	UnpackJSON([]byte(pack.Body.Data), &data)
	public := ParsePublic(string(Base64Decode(data.Public)))

	hash := pack.From.Sender.Hashname
	client.Connections[hash] = &Connect{
		connected: true,
		Chans: Chans{
			Action: make(chan bool),
			action: make(chan bool),
		},
		hashname:    hash,
		address:     pack.From.Address,
		public:      public,
		certificate: Base64Decode(data.Certificate),
		session:     DecryptRSA(client.keys.private, Base64Decode(data.Session)),
	}

	if pack.From.Hashname == pack.From.Sender.Hashname {
		client.Connections[hash].throwClient = public
		client.Connections[hash].relation = conn
	} else {
		client.Connections[hash].throwClient = client.Connections[pack.From.Hashname].public
	}
}

// Disconnect by GET option.
func (client *Client) disconnectGet(pack *Package) {
	hash := pack.From.Sender.Hashname
	if client.Connections[hash].relation != nil {
		client.Connections[hash].relation.Close()
	}
	delete(client.Connections, hash)
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

func (client *Client) signPackage(pack *Package) *Package {
	hash := HashSum(bytes.Join(
		[][]byte{
			[]byte(pack.Info.Network),
			[]byte(pack.Info.Version),
			[]byte(pack.From.Sender.Hashname),
			[]byte(pack.To.Receiver.Hashname),
			[]byte(pack.Head.Title),
			[]byte(pack.Head.Option),
			[]byte(pack.Body.Data),
			ToBytes(pack.Body.Desc.Id),
			[]byte(pack.Body.Desc.Rand),
		},
		[]byte{},
	))

	pack.Body.Test.Hash = Base64Encode(hash)
	pack.Body.Test.Sign = Base64Encode(Sign(client.keys.private, hash))
	return pack
}

// Encrypt package by session key. Encrypted data:
// 1) Head.Title;
// 2) Head.Option;
// 3) Body.Data;
// 4) Body.Desc.Rand;
func (client *Client) encryptPackage(pack *Package) *Package {
	var session []byte

	switch {
	case client.isConnected(pack.To.Receiver.Hashname):
		session = client.Connections[pack.To.Receiver.Hashname].session
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
				Id:          pack.Body.Desc.Id,
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
// 4) Body.Desc.Rand;
func (client *Client) decryptPackage(pack *Package) *Package {
	if !client.InConnections(pack.From.Sender.Hashname) {
		return nil
	}
	session := client.Connections[pack.From.Sender.Hashname].session
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
				Id:          pack.Body.Desc.Id,
				Rand:        string(DecryptAES(session, Base64Decode(pack.Body.Desc.Rand))),
				Hash:        pack.Body.Desc.Hash,
				Sign:        pack.Body.Desc.Sign,
				Nonce:       pack.Body.Desc.Nonce,
				Difficulty:  settings.DIFFICULTY,
				Redirection: pack.Body.Desc.Redirection,
			},
			Test: Test{
				Hash: pack.Body.Test.Hash,
				Sign: pack.Body.Test.Sign,
			},
		},
	}
}

// Get previous hash, generate random bytes, calculate new hash, sign and nonce for package.
// Save current hash in local storage.
func (client *Client) confirmPackage(random []byte, pack *Package) *Package {
	pack.Body.Desc.Id = client.Connections[pack.To.Receiver.Hashname].packageId
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
			ToBytes(pack.Body.Desc.Id),
			[]byte(pack.Body.Desc.Rand),
		},
		[]byte{},
	))
	pack.Body.Desc.Hash = Base64Encode(hash)
	pack.Body.Desc.Sign = Base64Encode(Sign(client.keys.private, hash))
	pack.Body.Desc.Nonce = ProofOfWork(hash, pack.Body.Desc.Difficulty)
	return pack
}

// Append information about network name, version.
// Append sender information: hashname, public, address.
func (client *Client) appendHeaders(pack *Package) *Package {
	pack.Info.Network = settings.NETWORK
	pack.Info.Version = settings.VERSION
	pack.From.Hashname = client.hashname
	pack.From.Address = client.address
	if pack.From.Sender.Hashname == "" {
		pack.From.Sender.Hashname = client.hashname
	}
	return pack
}

// Check if user connected to client.
func (client *Client) isConnected(hash string) bool {
	if _, ok := client.Connections[hash]; ok {
		return client.Connections[hash].connected
	}
	return false
}

// Get connection and check package.
func runServer(listener *Listener, handle func(*Client, *Package)) {
	defer listener.Close()
	for {
		conn, err := listener.listen.Accept()
		if err != nil {
			break
		}
		go serveNode(listener, handle, conn)
	}
}

// Read package by node.
func serveNode(listener *Listener, handle func(*Client, *Package), conn net.Conn) {
	var (
		client *Client
		hash   string
	)
	defer func() {
		if client != nil {
			delete(client.Connections, hash)
		}
		conn.Close()
	}()
	for {
		pack := readPackage(conn)
		if pack == nil {
			break
		}
		received := pack.receive(listener, handle, conn)
		if hash == "" && received {
			client = listener.Clients[pack.To.Hashname]
			hash = pack.From.Hashname
		}
	}
}

// Read package by client.
func serveClient(listener *Listener, client *Client, handle func(*Client, *Package), hash string, conn net.Conn) {
	defer func() {
		delete(client.Connections, hash)
		conn.Close()
	}()
	for {
		pack := readPackage(conn)
		if pack == nil {
			break
		}
		pack.receive(listener, handle, conn)
	}
}

// Read raw data and convert to package.
func readPackage(conn net.Conn) *Package {
	var (
		message string
		pack    = new(Package)
		size    = uint64(0)
		buffer  = make([]byte, settings.BUFF_SIZE)
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			return nil
		}
		size += uint64(length)
		if size >= settings.PACK_SIZE {
			return nil
		}
		message += string(buffer[:length])
		if strings.Contains(message, settings.END_BYTES) {
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

func readFile(filename string, id uint64) []byte {
	const BEGGINING = 0
	var FILE_SIZE = settings.PACK_SIZE / 4

	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	_, err = file.Seek(int64(id*FILE_SIZE), BEGGINING)
	if err != nil {
		return nil
	}

	var buffer = make([]byte, FILE_SIZE)
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
func printJSON(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(jsonData))
}

// For blockcipher encryption.
func paddingPKCS5(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

// For blockcipher decryption.
func unpaddingPKCS5(origData []byte) []byte {
    length := len(origData)
    if length == 0 {
        return nil
    }
    unpadding := int(origData[length-1])
    if length < unpadding {
        return nil
    }
    return origData[:(length - unpadding)]
}
