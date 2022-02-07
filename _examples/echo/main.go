package main

import (
	"fmt"
	"time"

	cr "github.com/number571/go-peer/crypto"
	lc "github.com/number571/go-peer/local"
	nt "github.com/number571/go-peer/network"
	gp "github.com/number571/go-peer/settings"
)

const (
	NODE_ADDRESS = ":8080"
)

var (
	ROUTE_MSG = []byte("/msg")
)

func main() {
	client := newNode()
	node := newNode()

	// Run node.
	go node.Listen(NODE_ADDRESS)
	time.Sleep(500 * time.Millisecond)

	// Connect to node.
	client.Connect(NODE_ADDRESS)

	// Create message and route.
	route := lc.NewRoute(node.Client().PubKey(), nil, nil)

	msg := lc.NewMessage(
		ROUTE_MSG,
		[]byte("hello, world!"),
	)

	// Send request 'hello, world!' to node.
	res, err := client.Broadcast(route, msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print response.
	fmt.Println(string(res))
}

func newNode() nt.Node {
	// Generate private key.
	settings := gp.NewSettings()

	priv := cr.NewPrivKey(settings.Get(gp.SizeAkey))
	node := nt.NewNode(lc.NewClient(priv))

	// Set local route to function.
	node.Handle(ROUTE_MSG, getMessage)
	return node
}

func getMessage(client lc.Client, msg lc.Message) []byte {
	// Receive message.
	hash := cr.LoadPubKey(msg.Head.Sender).Address()
	fmt.Printf("[%s] => '%s'\n", hash, msg.Body.Data)

	// Response.
	return msg.Body.Data
}
