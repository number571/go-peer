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
	DIFF_PACK = gp.Get("POWS_DIFF").(uint)
	ROUTE_MSG = []byte("/msg")
)

func main() {
	client := newClient()
	node := newClient()

	// Run node.
	go node.RunNode(NODE_ADDRESS)
	time.Sleep(500 * time.Millisecond)

	// Connect to node.
	client.Connect(NODE_ADDRESS)

	// Create message and route.
	msg := nt.NewMessage(
		ROUTE_MSG,
		[]byte("hello, world!"),
	).WithDiff(DIFF_PACK)
	route := nt.NewRoute(node.PubKey())

	// Send request 'hello, world!' to node.
	res, err := client.Send(msg, route)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print response.
	fmt.Println(string(res))
}

func newClient() *nt.Client {
	// Generate private key.
	priv := cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))
	client := nt.NewClient(priv)

	// Set local route to function.
	client.Handle(ROUTE_MSG, getMessage)
	return client
}

func getMessage(client *nt.Client, msg *nt.Message) []byte {
	// Receive message.
	hash := cr.LoadPubKey(msg.Head.Sender).Address()
	fmt.Printf("[%s] => '%s'\n", hash, msg.Body.Data)

	// Response.
	return msg.Body.Data
}
