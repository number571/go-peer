package gopeer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
)

// Get connection and check package.
func runServer(handle func(*Client, *Package), listener *Listener) {
	for {
		if listener.Setting.Listen == nil {
			break
		}
		conn, err := listener.Setting.Listen.Accept()
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
		wasEncrypted bool
	)

	pack = readPackage(conn)
	if pack == nil {
		return
	}

	client, ok := listener.Clients[pack.To.Receiver.Hashname]
	if !ok {
		return
	}

	pack, wasEncrypted = client.tryDecrypt(pack)
	// printJson(pack)
	if err := client.isValid(pack); err != nil {
		// fmt.Println(err)
		return
	}

	handleIsUsed := client.HandleAction(settings.TITLE_LASTHASH, pack,
		func(client *Client, pack *Package) (set string) {
			if !client.InConnections(pack.From.Sender.Hashname) {
				return
			}
			return client.Connections[pack.From.Sender.Hashname].LastHash
		},
		func(client *Client, pack *Package) {
			if !client.InConnections(pack.From.Sender.Hashname) {
				return
			}
			client.Connections[pack.From.Sender.Hashname].LastHash = pack.Body.Data
		},
	) 

	if handleIsUsed {
		return
	}

	handleIsUsed = client.HandleAction(settings.TITLE_CONNECT, pack,
		func(client *Client, pack *Package) (set string) {
			client.connectGet(pack)
			return set
		},
		func(client *Client, pack *Package) {
			if !client.InConnections(pack.From.Sender.Hashname) {
				return
			}
			client.Connections[pack.From.Sender.Hashname].Connected = true
			client.Connections[pack.From.Sender.Hashname].LastHash = pack.Body.Desc.CurrHash
		},
	)

	// Subsequent verification is carried out only if the data has been encrypted.
	if handleIsUsed || !wasEncrypted {
		return
	}

	client.Connections[pack.From.Sender.Hashname].LastHash = pack.Body.Desc.CurrHash

	switch pack.Head.Title {
	case settings.TITLE_DISCONNECT:
		switch pack.Head.Option {
		case settings.OPTION_GET:
			client.disconnectGet(pack)
			return
		case settings.OPTION_SET:
			client.Connections[pack.From.Sender.Hashname].Waiting <- true
			delete(client.Connections, pack.From.Sender.Hashname)
			return
		}
	}

	handle(client, pack)
}

// Connect by GET option.
func (client *Client) connectGet(pack *Package) {
	var (
		hash     = pack.From.Sender.Hashname
		lastHash = ""
	)
	if client.InConnections(hash) {
		lastHash = client.Connections[hash].LastHash
	}
	client.Connections[hash] = &Connect{
		Session:   DecryptRSA(client.Keys.Private, Base64Decode(pack.Body.Data)),
		Address:   pack.From.Address,
		LastHash:  lastHash,
		Connected: true,
		Public:    ParsePublic(string(Base64Decode(pack.From.Sender.Public))),
		Waiting:   make(chan bool),
	}
	client.Connections[pack.From.Sender.Hashname].LastHash = pack.Body.Desc.CurrHash
}

// Disconnect by GET option.
func (client *Client) disconnectGet(pack *Package) {
	client.Send(&Package{
		To: To{
			Receiver: Receiver{
				Hashname: pack.From.Sender.Hashname,
			},
			Address: pack.From.Address,
		},
		Head: Head{
			Title:  settings.TITLE_DISCONNECT,
			Option: settings.OPTION_SET,
		},
	})
	delete(client.Connections, pack.From.Sender.Hashname)
}

// Read raw data and convert to package.
func readPackage(conn net.Conn) *Package {
	var (
		message string
		pack    = new(Package)
		size 	= uint32(0)
		buffer  = make([]byte, settings.BUFFSIZE)
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			break
		}
		size += uint32(length)
		if size >= settings.PACKSIZE {
			return nil
		}
		message += string(buffer[:length])
	}
	err := json.Unmarshal(DecryptAES([]byte(settings.NOISE), []byte(message)), pack)
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
// 4) Body.Data.Desc;
func (client *Client) encryptPackage(pack *Package) *Package {
	var session []byte

	if pack.isLasthash() {
		return nil
	}

	switch {
	case client.IsConnected(pack.To.Receiver.Hashname):
		session = client.Connections[pack.To.Receiver.Hashname].Session
	case client.Connections[pack.To.Receiver.Hashname].PrevSession != nil:
		session = client.Connections[pack.To.Receiver.Hashname].PrevSession
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
				Public:   pack.From.Sender.Public,
			},
			Address: pack.From.Address,
		},
		To: To{
			Receiver: Receiver{
				Hashname: pack.To.Receiver.Hashname,
			},
			Address: pack.To.Address,
		},
		Head: Head{
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
	if DecryptAES(session, Base64Decode(pack.Head.Title)) == nil {
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
				Public:   pack.From.Sender.Public,
			},
			Address: pack.From.Address,
		},
		To: To{
			Receiver: Receiver{
				Hashname: pack.To.Receiver.Hashname,
			},
			Address: pack.To.Address,
		},
		Head: Head{
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

	public := ParsePublic(string(Base64Decode(pack.From.Sender.Public)))
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
		if pack.isLasthash() { // || pack.Head.Option == settings.OPTION_SET
			return nil
		}
		
		if client.Connections[pack.From.Sender.Hashname].LastHash != pack.Body.Desc.PrevHash {
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
	pack.From.Sender.Hashname = client.Hashname
	pack.From.Sender.Public = Base64Encode([]byte(StringPublic(client.Keys.Public)))
	pack.From.Address = client.Address
	return pack
}

// Get previous hash, generate random bytes, calculate new hash, sign and nonce for package.
// Save current hash in local storage.
func (client *Client) confirmPackage(pack *Package) *Package {
	if !pack.isLasthash() {
		pack.Body.Desc.PrevHash = client.Connections[pack.To.Receiver.Hashname].LastHash
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
		client.Connections[pack.To.Receiver.Hashname].LastHash = pack.Body.Desc.CurrHash
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
	client.Connections[hash].LastHash = ""
	client.Send(&Package{
		To: To{
			Receiver: Receiver{
				Hashname: hash,
			},
			Address: pack.To.Address,
		},
		Head: Head{
			Title:  settings.TITLE_LASTHASH,
			Option: settings.OPTION_GET,
		},
	})
	client.confirmPackage(client.appendHeaders(pack))
	var (
		savedPack = pack
	)
	if encPack := client.encryptPackage(pack); encPack != nil {
		pack = encPack
	}
	conn, err := net.Dial("tcp", pack.To.Address)
	if err != nil {
		delete(client.Connections, hash)
		return err
	}
	conn.Write(EncryptAES([]byte(settings.NOISE), PackJSON(pack)))
	conn.Close()
	if savedPack.Head.Option == settings.OPTION_GET {
		select {
		case <-client.Connections[hash].Waiting:
			err = nil
		case <-time.After(time.Duration(settings.RETRY_TIME) * time.Second):
			pack = savedPack
			if retryNum > 1 {
				retryNum--
				goto retry
			}
			err = fmt.Errorf("time is over (%d seconds)", settings.RETRY_TIME)
			delete(client.Connections, hash)
		}
	}
	return err
}

// For debug.
func printJson(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(jsonData))
}
