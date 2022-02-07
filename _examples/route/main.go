package main

import (
	"fmt"
	"sync"
	"time"

	cr "github.com/number571/go-peer/crypto"
	lc "github.com/number571/go-peer/local"
	nt "github.com/number571/go-peer/network"
	gp "github.com/number571/go-peer/settings"
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
	settings := gp.NewSettings()

	client1 := nt.NewNode(lc.NewClient(cr.NewPrivKey(settings.Get(gp.SizeAkey))))
	client2 := nt.NewNode(lc.NewClient(cr.NewPrivKey(settings.Get(gp.SizeAkey))))

	node1 := nt.NewNode(lc.NewClient(cr.NewPrivKey(settings.Get(gp.SizeAkey))))
	node2 := nt.NewNode(lc.NewClient(cr.NewPrivKey(settings.Get(gp.SizeAkey))))
	node3 := nt.NewNode(lc.NewClient(cr.NewPrivKey(settings.Get(gp.SizeAkey))))

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

	psender := cr.NewPrivKey(settings.Get(gp.SizeAkey))
	routes := []cr.PubKey{
		node1.Client().PubKey(),
		node2.Client().PubKey(),
		node3.Client().PubKey(),
	}

	var (
		wg    sync.WaitGroup
		count = 20
	)

	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int, sender, receiver nt.Node, psender cr.PrivKey, routes []cr.PubKey) {
			data := fmt.Sprintf("hello, world! [%d]", i)
			res, err := sender.Broadcast(
				lc.NewRoute(receiver.Client().PubKey(), psender, routes),
				lc.NewMessage(ROUTE_MSG, []byte(data)),
			)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(res))
			wg.Done()
		}(i, client1, client2, psender, routes)
	}
	wg.Wait()
}

func getMessage(client lc.Client, msg lc.Message) []byte {
	hash := cr.LoadPubKey(msg.Head.Sender).Address()
	fmt.Printf("[%s] => '%s'\n", hash, msg.Body.Data)
	return msg.Body.Data
}
