package gopeer

import (
	"net"
	"strings"
)

func NewNode(address string, client *Client) *Node {
	return &Node{
		address: address,
		client:  client,
	}
}

func NewPackage(title, data string) *Package {
	return &Package{
		Head: HeadPackage {
			Title: title,
		},
		Body: BodyPackage {
			Data: data,
		},
	}
}

func Handle(title string, client *Client, pack *Package, handle func(*Client, *Package) string) {
	switch pack.Head.Title {
	case title:
		client.send(client.Encrypt(
			ParsePublic(pack.Head.Sender), 
			NewPackage("_" + title, handle(client, pack)),
		))
	case "_" + title:
		client.response(
			ParsePublic(pack.Head.Sender),
			pack.Body.Data,
		)
	}
}

func (node *Node) Run() error {
	var err error
	node.listen, err = net.Listen("tcp", node.address)
	if err != nil {
		return err
	}
	defer node.listen.Close()
	for {
		conn, err := node.listen.Accept()
		if err != nil {
			break
		}
		if uint(len(node.client.connections)) > settings.CONN_SIZE {
			conn.Close()
			continue
		}
		node.client.connections[conn] = "client"
		go handleConn(conn, node.client, node.client.handle)
	}
	return nil
}

func handleConn(conn net.Conn, client *Client, handle func(*Client, *Package)) {
	defer func() {
		conn.Close()
		delete(client.connections, conn)
	}()
	for {
		pack := readPackage(conn)
		isRoute := false

checkAgain:

		if pack == nil {
			continue
		}

		client.mutex.Lock()
		if _, ok := client.mapping[pack.Body.Hash]; ok {
			client.mutex.Unlock()
			continue
		}
		if uint(len(client.mapping)) > settings.MAPP_SIZE {
			client.mapping = make(map[string]bool)
		}
		client.mapping[pack.Body.Hash] = true
		client.mutex.Unlock()

		if !ProofIsValid(Base64Decode(pack.Body.Hash), settings.POWS_DIFF, pack.Body.Npow) {
			continue
		}

		if isRoute {
			client.send(pack)
		} else {
			client.redirect(pack, conn)
		}
		
		decPack := client.Decrypt(pack)

		if decPack == nil {
			continue
		}

		client.mutex.Lock()
		if client.f2f.enabled && !client.InF2F(ParsePublic(decPack.Head.Sender)) {
			client.mutex.Unlock()
			continue
		}
		client.mutex.Unlock()

		if decPack.Head.Title == settings.ROUTE_MSG {
			pack = DeserializePackage(decPack.Body.Data)
			isRoute = true
			goto checkAgain
		}

		handle(client, decPack)
	}
}

func readPackage(conn net.Conn) *Package {
	var (
		message string
		size    = uint(0)
		buffer  = make([]byte, settings.BUFF_SIZE)
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			return nil
		}
		size += uint(length)
		if size > settings.PACK_SIZE {
			return nil
		}
		message += string(buffer[:length])
		if strings.Contains(message, settings.END_BYTES) {
			message = strings.Split(message, settings.END_BYTES)[0]
			break
		}
	}
	return DeserializePackage(message)
}
