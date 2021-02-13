package main

import (
	gp "./gopeer"
	"fmt"
	"time"
)

const (
	TITLE_MESSAGE = "MESSAGE"
	NODE1_ADDRESS = ":8080"
	NODE2_ADDRESS = ":9090"
)

/*
	client1 -> node2
	client1 <-> node1 <-> client2 <-> node2
*/

func main() {
	client1 := gp.NewClient(
		gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)),
		handleFunc,
	)
	client2 := gp.NewClient(
		gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)),
		handleFunc,
	)

	fmt.Println(gp.HashPublic(client1.Public()))

	node1 := gp.NewClient(
		gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)),
		handleFunc,
	)
	node2 := gp.NewClient(
		gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)),
		handleFunc,
	)

	go node1.RunNode(NODE1_ADDRESS)
	go node2.RunNode(NODE2_ADDRESS)
	time.Sleep(500 * time.Millisecond)

	client1.Connect(NODE1_ADDRESS)
	client2.Connect(NODE1_ADDRESS)

	client2.Connect(NODE2_ADDRESS)

	for i := 0; i < 10; i++ {
		res, err := client1.Send(
			node2.Public(), 
			gp.NewPackage(TITLE_MESSAGE, fmt.Sprintf("hello, world! [%d]", i)),
			nil, 
			nil, 
		)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(res)
	}
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	client.Handle(TITLE_MESSAGE, pack, getMessage)
}

func getMessage(client *gp.Client, pack *gp.Package) (set string) {
	public := gp.ParsePublic(pack.Head.Sender)
	fmt.Printf("[%s] => '%s'\n", gp.HashPublic(public), pack.Body.Data)
	return "ok"
}
