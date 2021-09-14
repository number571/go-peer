package main

import (
	"fmt"
	"time"

	gp "github.com/number571/gopeer"
	cr "github.com/number571/gopeer/crypto"
	nt "github.com/number571/gopeer/network"
)

const (
	NODE_ADDRESS = ":8080"
)

var (
	ROUTE_MSG = []byte("/msg")
)

func main() {
	client1 := nt.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))).Handle(ROUTE_MSG, getMessage)
	client2 := nt.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))).Handle(ROUTE_MSG, getMessage)
	clinode := nt.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))).Handle(ROUTE_MSG, getMessage)

	go clinode.RunNode(NODE_ADDRESS)
	time.Sleep(500 * time.Millisecond)

	client1.Connect(NODE_ADDRESS)
	client2.Connect(NODE_ADDRESS)

	diff := gp.Get("POWS_DIFF").(uint)
	res, err := client1.Send(
		nt.NewMessage(ROUTE_MSG, []byte("hello, world!")).WithDiff(diff),
		nt.NewRoute(client2.PubKey()),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(res))
}

func getMessage(client *nt.Client, msg *nt.Message) []byte {
	hash := cr.LoadPubKey(msg.Head.Sender).Address()
	fmt.Printf("[%s] => '%s'\n", hash, msg.Body.Data)
	return msg.Body.Data
}
