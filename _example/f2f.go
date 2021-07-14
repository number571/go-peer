package main

import (
	"fmt"
	"time"

	gp "./gopeer"
)

const (
	TITLE_MESSAGE = "TITLE_MESSAGE"
	NODE_ADDRESS  = ":8080"
)

func main() {
	client1 := gp.NewClient(gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)))
	client2 := gp.NewClient(gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)))
	clinode := gp.NewClient(gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)))

	fmt.Println(client1.F2F.State(), client2.F2F.State())

	client1.F2F.Switch()
	client2.F2F.Switch()

	fmt.Println(client1.F2F.State(), client2.F2F.State())

	client1.F2F.Append(client2.PublicKey())
	client2.F2F.Append(client1.PublicKey())

	client1.Handle(TITLE_MESSAGE, getMessage)
	client2.Handle(TITLE_MESSAGE, getMessage)
	clinode.Handle(TITLE_MESSAGE, getMessage)

	go clinode.RunNode(NODE_ADDRESS)

	time.Sleep(500 * time.Millisecond)

	client1.Connect(NODE_ADDRESS)
	client2.Connect(NODE_ADDRESS)

	res, err := client1.Send(
		client2.PublicKey(),
		gp.NewPackage(TITLE_MESSAGE, []byte("hello, world!")),
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
