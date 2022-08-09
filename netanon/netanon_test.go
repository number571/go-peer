package netanon

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/number571/go-peer/client"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/friends"
	"github.com/number571/go-peer/network"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/queue"
	"github.com/number571/go-peer/testutils"
)

const (
	tcWait = 30
	tcIter = 10
)

func TestComplex(t *testing.T) {
	nodes := testNewNodes(tcWait)
	defer testFreeNodes(nodes[:])

	wg := sync.WaitGroup{}
	wg.Add(tcIter)

	for i := 0; i < tcIter; i++ {
		go func(i int) {
			defer wg.Done()
			reqBody := fmt.Sprintf("%s (%d)", testutils.TcBody, i)

			// nodes[1] -> nodes[0] -> nodes[2]
			resp, err := nodes[0].Request(
				nodes[1].Client().PubKey(),
				payload.NewPayload(uint64(testutils.TcHead), []byte(reqBody)),
			)
			if err != nil {
				t.Errorf("%s (%d)", err.Error(), i)
				return
			}

			if string(resp) != fmt.Sprintf("%s (response)", reqBody) {
				t.Errorf("string(resp) != reqBody (%d)", i)
				return
			}
		}(i)
	}

	wg.Wait()
}

func TestF2F(t *testing.T) {
	nodes := testNewNodes(tcWait)
	defer testFreeNodes(nodes[:])

	nodes[0].F2F().Switch(true)
	nodes[1].F2F().Switch(true)

	nodes[0].F2F().Append(nodes[1].Client().PubKey())
	nodes[1].F2F().Append(nodes[0].Client().PubKey())

	testRequest(t, 1, nodes)
}

func TestF2FWithoutFriends(t *testing.T) {
	// 5 seconds for wait
	nodes := testNewNodes(5)
	defer testFreeNodes(nodes[:])

	nodes[0].F2F().Switch(true)
	nodes[1].F2F().Switch(true)

	testRequest(t, 2, nodes)
}

func testRequest(t *testing.T, mode int, nodes [5]INode) {
	reqBody := fmt.Sprintf("%s (%d)", testutils.TcBody, mode)

	// nodes[1] -> nodes[0] -> nodes[2]
	resp, err := nodes[0].Request(
		nodes[1].Client().PubKey(),
		payload.NewPayload(uint64(testutils.TcHead), []byte(reqBody)),
	)
	if err != nil {
		if mode == 2 {
			return
		}
		t.Errorf("%s (mode=%d)", err.Error(), mode)
		return
	}

	if string(resp) != fmt.Sprintf("%s (response)", reqBody) {
		t.Errorf("string(resp) != reqBody")
		return
	}
}

// nodes[0], nodex[1] = clients
// nodes[2], nodes[3], nodes[4] = routes
// (nodes[0]) -> nodes[2] -> nodes[3] -> nodes[4] -> (nodes[1])
func testNewNodes(secondsWait int) [5]INode {
	nodes := [5]INode{}

	clients := testNewClients()
	for i := 0; i < 5; i++ {
		nodes[i] = testNewNode(i, secondsWait, clients)
	}

	for _, node := range nodes {
		node.Handle(
			testutils.TcHead,
			func(node INode, sender asymmetric.IPubKey, pl payload.IPayload) []byte {
				// send response
				resp := fmt.Sprintf("%s (response)", string(pl.Body()))
				return []byte(resp)
			},
		)
	}

	go func() {
		err := nodes[2].Network().Listen(testutils.TgAddrs[2])
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		err := nodes[4].Network().Listen(testutils.TgAddrs[3])
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(200 * time.Millisecond)

	// nodes to routes
	nodes[0].Network().Connect(testutils.TgAddrs[2])
	nodes[1].Network().Connect(testutils.TgAddrs[3])

	// routes to routes
	nodes[3].Network().Connect(testutils.TgAddrs[2])
	nodes[3].Network().Connect(testutils.TgAddrs[3])

	return nodes
}

func testNewClients() [5]client.IClient {
	clients := [5]client.IClient{}
	for i := 0; i < 5; i++ {
		clients[i] = client.NewClient(
			client.NewSettings(10, (1<<10)),
			asymmetric.NewRSAPrivKey(1024),
		)
	}
	return clients
}

func testNewNode(i, secondsWait int, clients [5]client.IClient) INode {
	msgSize := uint64(1 << 20)
	return NewNode(
		NewSettings(
			1,
			3,
			time.Duration(secondsWait)*time.Second,
		),
		clients[i],
		network.NewNode(network.NewSettings(
			msgSize,
			3,
			1024,
			10,
			20,
			5*time.Second,
		)),
		queue.NewQueue(
			queue.NewSettings(
				20,
				10,
				msgSize,
				300*time.Millisecond,
			),
			clients[i],
		),
		friends.NewF2F(),
		func() []asymmetric.IPubKey {
			return testGetPubKeys(clients[2:])
		},
	)
}

func testGetPubKeys(clients []client.IClient) []asymmetric.IPubKey {
	pubKeys := []asymmetric.IPubKey{}
	for _, client := range clients {
		pubKeys = append(pubKeys, client.PubKey())
	}
	return pubKeys
}

func testFreeNodes(nodes []INode) {
	for _, node := range nodes {
		node.Close()
	}
}
