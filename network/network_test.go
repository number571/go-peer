package network

import (
	"sync"
	"testing"
	"time"

	cr "github.com/number571/go-peer/crypto"
	lc "github.com/number571/go-peer/local"
	gp "github.com/number571/go-peer/settings"
	tu "github.com/number571/go-peer/settings/testutils"
)

const (
	nodeAddress1 = ":7070"
	nodeAddress2 = ":8080"
	nodeAddress3 = ":9090"
)

var (
	routeEcho = []byte("/echo")
)

func echoMessage(client lc.Client, msg lc.Message) []byte {
	return msg.Body.Data
}

func newNode() Node {
	settings := tu.NewSettings()
	privKey := cr.NewPrivKey(1024)
	client := lc.NewClient(privKey, settings)
	return NewNode(client)
}

// Simple broadcast testing

func initSimple() ([3]Node, lc.Route, lc.Message) {
	client1 := newNode()
	client2 := newNode()

	node1 := newNode()

	client1.Handle(routeEcho, echoMessage)
	client2.Handle(routeEcho, echoMessage)

	go node1.Listen(nodeAddress1)

	time.Sleep(500 * time.Millisecond)

	client1.Connect(nodeAddress1)
	client2.Connect(nodeAddress1)

	return [3]Node{client1, client2, node1},
		lc.NewRoute(client2.Client().PubKey(), nil, nil),
		lc.NewMessage(routeEcho, []byte("hello, world!"))
}

func TestSimple(t *testing.T) {
	nodes, route, msg := initSimple()
	defer nodes[2].Close()

	_, err := nodes[0].Broadcast(route, msg)
	if err != nil {
		t.Error(err)
		return
	}
}

func BenchmarkSimple(b *testing.B) {
	var wg sync.WaitGroup

	nodes, route, msg := initSimple()
	defer nodes[2].Close()

	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func(i int, sender Node, route lc.Route, msg lc.Message) {
			_, err := sender.Broadcast(route, msg)
			if err != nil {
				b.Error(err)
				return
			}
			wg.Done()
		}(i, nodes[0], route, msg)
	}
	wg.Wait()
}

// F2F testing

func TestF2F(t *testing.T) {
	nodes, route, msg := initSimple()
	defer nodes[2].Close()

	// time wait = 1 second
	nodes[0].Client().Settings().Set(gp.TimeWait, 1)
	nodes[1].Client().Settings().Set(gp.TimeWait, 1)

	nodes[0].F2F().Switch()
	nodes[1].F2F().Switch()

	_, err := nodes[0].Broadcast(route, msg)
	if err == nil {
		t.Errorf("f2f mode not working")
		return
	}

	nodes[0].F2F().Append(nodes[1].Client().PubKey())
	nodes[1].F2F().Append(nodes[0].Client().PubKey())

	_, err = nodes[0].Broadcast(route, msg)
	if err != nil {
		t.Error(err)
		return
	}
}

// Route testing

func initRoute() ([5]Node, lc.Route, lc.Message) {
	client1 := newNode()
	client2 := newNode()

	node1 := newNode()
	node2 := newNode()
	node3 := newNode()

	client1.Handle(routeEcho, echoMessage)
	client2.Handle(routeEcho, echoMessage)

	go node1.Listen(nodeAddress1)
	go node2.Listen(nodeAddress2)
	go node3.Listen(nodeAddress3)

	time.Sleep(500 * time.Millisecond)

	node1.Connect(nodeAddress2)
	node2.Connect(nodeAddress3)

	client1.Connect(nodeAddress1)
	client2.Connect(nodeAddress3)

	psender := cr.NewPrivKey(client1.Client().PubKey().Size())
	routes := []cr.PubKey{
		node1.Client().PubKey(),
		node2.Client().PubKey(),
		node3.Client().PubKey(),
	}

	return [5]Node{client1, client2, node1, node2, node3},
		lc.NewRoute(client2.Client().PubKey(), psender, routes),
		lc.NewMessage(routeEcho, []byte("hello, world!"))
}

func TestRoute(t *testing.T) {
	nodes, route, msg := initRoute()
	defer nodes[2].Close()
	defer nodes[3].Close()
	defer nodes[4].Close()

	_, err := nodes[0].Broadcast(route, msg)
	if err != nil {
		t.Error(err)
		return
	}
}

func BenchmarkRoute(b *testing.B) {
	var wg sync.WaitGroup

	nodes, route, msg := initRoute()
	defer nodes[2].Close()
	defer nodes[3].Close()
	defer nodes[4].Close()

	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func(i int, sender Node, route lc.Route, msg lc.Message) {
			_, err := sender.Broadcast(route, msg)
			if err != nil {
				b.Error(err)
				return
			}
			wg.Done()
		}(i, nodes[0], route, msg)
	}
	wg.Wait()
}
