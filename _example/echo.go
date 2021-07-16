package main

import (
	"fmt"
	"time"

	gp "github.com/number571/gopeer"
)

const (
	ROUTE_MSG    = "/msg"
	NODE_ADDRESS = ":8080"
)

func main() {
	client1 := gp.NewClient(gp.GenerateKey(gp.Get("AKEY_SIZE").(uint))).Handle(ROUTE_MSG, getMessage)
	client2 := gp.NewClient(gp.GenerateKey(gp.Get("AKEY_SIZE").(uint))).Handle(ROUTE_MSG, getMessage)
	clinode := gp.NewClient(gp.GenerateKey(gp.Get("AKEY_SIZE").(uint))).Handle(ROUTE_MSG, getMessage)

	go clinode.RunNode(NODE_ADDRESS)
	time.Sleep(500 * time.Millisecond)

	client1.Connect(NODE_ADDRESS)
	client2.Connect(NODE_ADDRESS)

	res, err := client1.Send(
		client2.PublicKey(),
		gp.NewPackage(ROUTE_MSG, []byte("hello, world!")),
		nil,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(res))
}

func getMessage(client *gp.Client, pack *gp.Package) []byte {
	hash := gp.HashPublicKey(gp.BytesToPublicKey(pack.Head.Sender))
	fmt.Printf("[%s] => '%s'\n", hash, pack.Body.Data)
	return pack.Body.Data
}
