package main

import (
	"fmt"
	"time"

	gp "github.com/number571/gopeer"
	cr "github.com/number571/gopeer/crypto"
)

const (
	NODE_ADDRESS = ":8080"
)

var (
	ROUTE_MSG = []byte("/msg")
)

func main() {
	client1 := gp.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint)))
	client2 := gp.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint)))
	clinode := gp.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint)))

	fmt.Println(client1.F2F.State(), client2.F2F.State())

	client1.F2F.Switch()
	client2.F2F.Switch()

	fmt.Println(client1.F2F.State(), client2.F2F.State())

	client1.F2F.Append(client2.PubKey())
	client2.F2F.Append(client1.PubKey())

	client1.Handle(ROUTE_MSG, getMessage)
	client2.Handle(ROUTE_MSG, getMessage)
	clinode.Handle(ROUTE_MSG, getMessage)

	go clinode.RunNode(NODE_ADDRESS)

	time.Sleep(500 * time.Millisecond)

	client1.Connect(NODE_ADDRESS)
	client2.Connect(NODE_ADDRESS)

	res, err := client1.Send(
		gp.NewPackage(ROUTE_MSG, []byte("hello, world!")),
		gp.NewRoute(client2.PubKey()),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(res))
}

func getMessage(client *gp.Client, pack *gp.Package) []byte {
	hash := cr.LoadPubKey(pack.Head.Sender).Address()
	fmt.Printf("[%s] => '%s'\n", hash, pack.Body.Data)
	return pack.Body.Data
}
