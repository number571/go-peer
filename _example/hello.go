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
	client := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))
	node := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))

	go gp.NewListener(NODE_ADDRESS, node).Run(handleFunc)
	time.Sleep(500 * time.Millisecond)

	client.Connect(NODE_ADDRESS, handleFunc)

	res, err := client.Send(node.Public(), &gp.Package{
		Head: gp.HeadPackage{
			Title: TITLE_MESSAGE,
		},
		Body: gp.BodyPackage{
			Data: "hello, world!",
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	gp.Handle(TITLE_MESSAGE, client, pack, getMessage)
}

func getMessage(client *gp.Client, pack *gp.Package) (set string) {
	public := gp.ParsePublic(pack.Head.Sender)
	fmt.Printf("[%s] => '%s'\n", gp.HashPublic(public), pack.Body.Data)
	return "ok"
}
