package main

import (
	gp "./gopeer"
	"fmt"
	"time"
)

const (
	TITLE_MESSAGE = "TITLE_MESSAGE"
	NODE_ADDRESS  = ":8080"
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

	fmt.Println(client1.F2F.State(), client2.F2F.State())

	client1.F2F.Switch()
	client2.F2F.Switch()

	fmt.Println(client1.F2F.State(), client2.F2F.State())

	client1.F2F.Append(client2.PublicKey())
	client2.F2F.Append(client1.PublicKey())

	node := gp.NewClient(
		gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)),
		nil,
	)
	go node.RunNode(NODE_ADDRESS)

	time.Sleep(500 * time.Millisecond)

	client1.Connect(NODE_ADDRESS)
	client2.Connect(NODE_ADDRESS)

	res, err := client1.Send(
		client2.PublicKey(),
		gp.NewPackage(TITLE_MESSAGE, "hello, world!"),
		nil,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	client.Handle(TITLE_MESSAGE, pack, getMessage)
}

func getMessage(client *gp.Client, pack *gp.Package) (set string) {
	hash := gp.HashPublicKey(gp.StringToPublicKey(pack.Head.Sender))
	fmt.Printf("[%s] => '%s'\n", hash, pack.Body.Data)
	return pack.Body.Data
}
