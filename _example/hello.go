package main

import (
	gp "./gopeer"
	"fmt"
	"time"
)

const (
	GET_MESSAGE  = "GET_MESSAGE"
	SET_MESSAGE  = "SET_MESSAGE"
	NODE_ADDRESS = ":8080"
)

func main() {
	client := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))
	node := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))

	go gp.NewListener(NODE_ADDRESS, node).Run(handleFunc)
	time.Sleep(500 * time.Millisecond)

	client.Connect(NODE_ADDRESS, handleFunc)

	err := client.Request(node.Public(), &gp.Package{
		Head: gp.HeadPackage{
			Title: GET_MESSAGE,
		},
		Body: gp.BodyPackage{
			Data: "hello, world!",
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	switch pack.Head.Title {
	case GET_MESSAGE:
		public := gp.ParsePublic(pack.Head.Sender)
		fmt.Printf("[%s] => '%s'\n", gp.HashPublic(public), pack.Body.Data)
		client.Send(public, &gp.Package{
			Head: gp.HeadPackage{
				Title: SET_MESSAGE,
			},
		})
	case SET_MESSAGE:
		client.Response(gp.ParsePublic(pack.Head.Sender))
	default:
		fmt.Println("title undefined")
	}
}
