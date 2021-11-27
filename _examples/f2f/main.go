package main

import (
	"fmt"
	"time"

	gp "github.com/number571/gopeer"
	cr "github.com/number571/gopeer/crypto"
	lc "github.com/number571/gopeer/local"
	nt "github.com/number571/gopeer/network"
)

const (
	NODE_ADDRESS = ":8080"
)

var (
	ROUTE_MSG = []byte("/msg")
)

func main() {
	node1 := nt.NewNode(lc.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))))
	node2 := nt.NewNode(lc.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))))
	lnode := nt.NewNode(lc.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))))

	fmt.Println(node1.F2F().State(), node2.F2F().State())

	node1.F2F().Switch()
	node2.F2F().Switch()

	fmt.Println(node1.F2F().State(), node2.F2F().State())

	node1.F2F().Append(node2.Client().PubKey())
	node2.F2F().Append(node1.Client().PubKey())

	node1.Handle(ROUTE_MSG, getMessage)
	node2.Handle(ROUTE_MSG, getMessage)
	lnode.Handle(ROUTE_MSG, getMessage)

	go lnode.Listen(NODE_ADDRESS)

	time.Sleep(500 * time.Millisecond)

	node1.Connect(NODE_ADDRESS)
	node2.Connect(NODE_ADDRESS)

	diff := gp.Get("POWS_DIFF").(uint)
	res, err := node1.Send(
		lc.NewMessage(ROUTE_MSG, []byte("hello, world!"), diff),
		lc.NewRoute(node2.Client().PubKey()),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(res))
}

func getMessage(client *lc.Client, msg *lc.Message) []byte {
	hash := cr.LoadPubKey(msg.Head.Sender).Address()
	fmt.Printf("[%s] => '%s'\n", hash, msg.Body.Data)
	return msg.Body.Data
}
