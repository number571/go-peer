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

	node := gp.NewClient(
		gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)),
		handleFunc,
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
	publicBytes := gp.Base64Decode(pack.Head.Sender)
	hash := gp.Base64Encode(gp.HashSum(publicBytes))
	fmt.Printf("[%s] => '%s'\n", hash, pack.Body.Data)
	return pack.Body.Data
}
