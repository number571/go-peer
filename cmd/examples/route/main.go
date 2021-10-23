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
	NODE1_ADDRESS = ":7070"
	NODE2_ADDRESS = ":8080"
	NODE3_ADDRESS = ":9090"
)

var (
	ROUTE_MSG = []byte("/msg")
)

func main() {
	client1 := nt.NewNode(lc.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))))
	client2 := nt.NewNode(lc.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))))

	node1 := nt.NewNode(lc.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))))
	node2 := nt.NewNode(lc.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))))
	node3 := nt.NewNode(lc.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))))

	client1.Handle(ROUTE_MSG, getMessage)
	client2.Handle(ROUTE_MSG, getMessage)

	node1.Handle(ROUTE_MSG, getMessage)
	node2.Handle(ROUTE_MSG, getMessage)
	node3.Handle(ROUTE_MSG, getMessage)

	go node1.Listen(NODE1_ADDRESS)
	go node2.Listen(NODE2_ADDRESS)
	go node3.Listen(NODE3_ADDRESS)

	time.Sleep(500 * time.Millisecond)

	node1.Connect(NODE2_ADDRESS)
	node2.Connect(NODE3_ADDRESS)

	client1.Connect(NODE1_ADDRESS)
	client2.Connect(NODE3_ADDRESS)

	psender := cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))
	routes := []cr.PubKey{
		node1.Client().PubKey(),
		node2.Client().PubKey(),
		node3.Client().PubKey(),
	}

	diff := gp.Get("POWS_DIFF").(uint)
	for i := 0; i < 10; i++ {
		res, err := client1.Send(
			lc.NewMessage(ROUTE_MSG, []byte("hello, world!"), diff),
			lc.NewRoute(client2.Client().PubKey()).WithRoad(psender, routes),
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(res))
	}
}

func getMessage(client *lc.Client, msg *lc.Message) []byte {
	hash := cr.LoadPubKey(msg.Head.Sender).Address()
	fmt.Printf("[%s] => '%s'\n", hash, msg.Body.Data)
	return msg.Body.Data
}
