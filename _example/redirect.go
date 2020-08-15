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
	client1 := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))
	client2 := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))

	fmt.Println(gp.HashPublic(client1.Public()))

	node1 := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))
	node2 := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))

	go gp.NewListener(NODE1_ADDRESS, node1).Run(handleFunc)
	go gp.NewListener(NODE2_ADDRESS, node2).Run(handleFunc)
	time.Sleep(500 * time.Millisecond)

	client1.Connect(NODE1_ADDRESS, handleFunc)
	client2.Connect(NODE1_ADDRESS, handleFunc)

	client2.Connect(NODE2_ADDRESS, handleFunc)

	for i := 0; i < 10; i++ {
		res, err := client1.Send(node2.Public(), &gp.Package{
			Head: gp.HeadPackage{
				Title: TITLE_MESSAGE,
			},
			Body: gp.BodyPackage{
				Data: fmt.Sprintf("hello, world! [%d]", i),
			},
		})
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(res)
	}
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	gp.Handle(TITLE_MESSAGE, client, pack, getMessage)
}

func getMessage(client *gp.Client, pack *gp.Package) (set string) {
	public := gp.ParsePublic(pack.Head.Sender)
	fmt.Printf("[%s] => '%s'\n", gp.HashPublic(public), pack.Body.Data)
	return "ok"
}
