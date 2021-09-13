package main

import (
	"crypto/rsa"
	"fmt"
	"time"

	gp "github.com/number571/gopeer"
)

const (
	TITLE_MESSAGE = "TITLE_MESSAGE"
	NODE1_ADDRESS = ":7070"
	NODE2_ADDRESS = ":8080"
	NODE3_ADDRESS = ":9090"
)

func main() {
	client1 := gp.NewClient(gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)))
	client2 := gp.NewClient(gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)))

	node1 := gp.NewClient(gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)))
	node2 := gp.NewClient(gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)))
	node3 := gp.NewClient(gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)))

	client1.Handle(TITLE_MESSAGE, getMessage)
	client2.Handle(TITLE_MESSAGE, getMessage)

	node1.Handle(TITLE_MESSAGE, getMessage)
	node2.Handle(TITLE_MESSAGE, getMessage)
	node3.Handle(TITLE_MESSAGE, getMessage)

	go node1.RunNode(NODE1_ADDRESS)
	go node2.RunNode(NODE2_ADDRESS)
	go node3.RunNode(NODE3_ADDRESS)

	time.Sleep(500 * time.Millisecond)

	node1.Connect(NODE2_ADDRESS)
	node2.Connect(NODE3_ADDRESS)

	client1.Connect(NODE1_ADDRESS)
	client2.Connect(NODE3_ADDRESS)

	psender := gp.GenerateKey(gp.Get("AKEY_SIZE").(uint))
	routes := []*rsa.PublicKey{
		node1.PublicKey(),
		node2.PublicKey(),
		node3.PublicKey(),
	}

	for i := 0; i < 10; i++ {
		res, err := client1.Send(
			gp.NewPackage(TITLE_MESSAGE, []byte("hello, world!")),
			gp.NewRoute(client2.PublicKey()).Psender(psender).Routes(routes),
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(res))
	}
}

func getMessage(client *gp.Client, pack *gp.Package) []byte {
	hash := gp.HashPublicKey(gp.BytesToPublicKey(pack.Head.Sender))
	fmt.Printf("[%s] => '%s'\n", hash, pack.Body.Data)
	return pack.Body.Data
}
