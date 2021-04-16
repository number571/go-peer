package main

import (
	gp "./gopeer"
	"crypto/rsa"
	"fmt"
	"time"
)

const (
	TITLE_MESSAGE = "TITLE_MESSAGE"
	NODE1_ADDRESS = ":7070"
	NODE2_ADDRESS = ":8080"
	NODE3_ADDRESS = ":9090"
)

func main() {
	client1 := gp.NewClient(
		gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)),
		handleFunc,
	)
	client2 := gp.NewClient(
		gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)),
		handleFunc,
	)

	node1 := gp.NewClient(
		gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)),
		handleFunc,
	)
	go node1.RunNode(NODE1_ADDRESS)

	node2 := gp.NewClient(
		gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)),
		handleFunc,
	)
	go node2.RunNode(NODE2_ADDRESS)

	node3 := gp.NewClient(
		gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)),
		handleFunc,
	)
	go node3.RunNode(NODE3_ADDRESS)
	time.Sleep(500 * time.Millisecond)

	node1.Connect(NODE2_ADDRESS)
	node2.Connect(NODE3_ADDRESS)

	client1.Connect(NODE1_ADDRESS)
	client2.Connect(NODE3_ADDRESS)

	pseudoSender := gp.GenerateKey(gp.Get("AKEY_SIZE").(uint))
	route := []*rsa.PublicKey{
		node1.PublicKey(),
		node2.PublicKey(),
		node3.PublicKey(),
	}

	for i := 0; i < 10; i++ {
		res, err := client1.Send(
			client2.PublicKey(),
			gp.NewPackage(TITLE_MESSAGE, []byte("hello, world!")),
			route,
			pseudoSender,
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(res))
	}
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	client.Handle(TITLE_MESSAGE, pack, getMessage)
}

func getMessage(client *gp.Client, pack *gp.Package) []byte {
	hash := gp.HashPublicKey(gp.BytesToPublicKey(pack.Head.Sender))
	fmt.Printf("[%s] => '%s'\n", hash, pack.Body.Data)
	return pack.Body.Data
}
