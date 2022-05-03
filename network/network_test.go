package network

import (
	"sync"
	"testing"
	"time"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/settings"
	"github.com/number571/go-peer/settings/testutils"
)

const (
	tcNodeAddress1 = ":7070"
	tcNodeAddress2 = ":8080"
	tcNodeAddress3 = ":9090"
)

var (
	tgRouteEcho = []byte("/echo")
)

func echoMessage(node INode, msg local.IMessage) []byte {
	return msg.Body().Data()
}

func newNode() INode {
	sett := testutils.NewSettings()
	privKey := crypto.NewPrivKey(1024)
	client := local.NewClient(privKey, sett)
	return NewNode(client)
}

// Simple broadcast testing

func initSimple() ([3]INode, local.IRoute, local.IMessage) {
	client1 := newNode()
	client2 := newNode()

	node1 := newNode()

	client1.Handle(tgRouteEcho, echoMessage)
	client2.Handle(tgRouteEcho, echoMessage)

	go node1.Listen(tcNodeAddress1)

	time.Sleep(200 * time.Millisecond)

	client1.Connect(tcNodeAddress1)
	client2.Connect(tcNodeAddress1)

	return [3]INode{client1, client2, node1},
		local.NewRoute(client2.Client().PubKey()),
		local.NewMessage(tgRouteEcho, []byte("hello, world!"))
}

func TestSimple(t *testing.T) {
	nodes, route, msg := initSimple()
	defer nodes[2].Close()

	_, err := nodes[0].Request(route, msg)
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
		go func(i int, sender INode, route local.IRoute, msg local.IMessage) {
			_, err := sender.Request(route, msg)
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

	// time wait = timePseudo+1 second
	timeOut := nodes[0].Client().Settings().Get(settings.TimePrsp) + 1

	nodes[0].Client().Settings().Set(settings.TimeWait, timeOut)
	nodes[1].Client().Settings().Set(settings.TimeWait, timeOut)

	nodes[0].F2F().Switch(true)
	nodes[1].F2F().Switch(true)

	_, err := nodes[0].Request(route, msg)
	if err == nil {
		t.Errorf("f2f mode not working")
		return
	}

	nodes[0].F2F().Append(nodes[1].Client().PubKey())
	nodes[1].F2F().Append(nodes[0].Client().PubKey())

	_, err = nodes[0].Request(route, msg)
	if err != nil {
		t.Error(err)
		return
	}
}

// TEST CHECKER

func TestChecker(t *testing.T) {
	nodes, _, _ := initSimple()
	defer nodes[2].Close()

	nodes[1].Online().Switch(true)

	nodes[0].Checker().Switch(true)
	nodes[0].Checker().Append(nodes[1].Client().PubKey())

	// sleep = timePseudo+1 second
	timeOut := nodes[0].Client().Settings().Get(settings.TimePrsp) + 1
	time.Sleep(time.Duration(timeOut) * time.Second)

	info := nodes[0].Checker().ListWithInfo()[0]
	if !info.Online() {
		t.Errorf("checker not working")
	}
}

// TEST PSEUDO

func TestPseudo(t *testing.T) {
	nodes, route, msg := initSimple()
	defer nodes[2].Close()

	for _, node := range nodes {
		node.Pseudo().Switch(true)
	}

	for i := 0; i < 3; i++ {
		_, err := nodes[0].Request(route, msg)
		if err != nil {
			t.Error(err)
			return
		}
	}
}

// TEST ROUTE

func initRoute() ([5]INode, local.IRoute, local.IMessage) {
	client1 := newNode()
	client2 := newNode()

	node1 := newNode()
	node2 := newNode()
	node3 := newNode()

	client1.Handle(tgRouteEcho, echoMessage)
	client2.Handle(tgRouteEcho, echoMessage)

	go node1.Listen(tcNodeAddress1)
	go node2.Listen(tcNodeAddress2)
	go node3.Listen(tcNodeAddress3)

	time.Sleep(200 * time.Millisecond)

	node1.Connect(tcNodeAddress2)
	node2.Connect(tcNodeAddress3)

	client1.Connect(tcNodeAddress1)
	client2.Connect(tcNodeAddress3)

	client2.WithResponseRouter(func(_ INode) []crypto.IPubKey {
		return []crypto.IPubKey{
			node2.Client().PubKey(),
			node3.Client().PubKey(),
			node1.Client().PubKey(),
		}
	})

	psender := crypto.NewPrivKey(client1.Client().PubKey().Size())
	routes := []crypto.IPubKey{
		node1.Client().PubKey(),
		node2.Client().PubKey(),
		node3.Client().PubKey(),
	}

	return [5]INode{client1, client2, node1, node2, node3},
		local.NewRoute(client2.Client().PubKey()).WithRedirects(psender, routes),
		local.NewMessage(tgRouteEcho, []byte("hello, world!"))
}

func TestRoute(t *testing.T) {
	nodes, route, msg := initRoute()
	defer nodes[2].Close()
	defer nodes[3].Close()
	defer nodes[4].Close()

	_, err := nodes[0].Request(route, msg)
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
		go func(i int, sender INode, route local.IRoute, msg local.IMessage) {
			_, err := sender.Request(route, msg)
			if err != nil {
				b.Error(err)
				return
			}
			wg.Done()
		}(i, nodes[0], route, msg)
	}
	wg.Wait()
}
