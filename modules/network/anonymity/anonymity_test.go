package anonymity

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/client"
	"github.com/number571/go-peer/modules/closer"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/friends"
	"github.com/number571/go-peer/modules/network"
	"github.com/number571/go-peer/modules/payload"
	"github.com/number571/go-peer/modules/queue"
	"github.com/number571/go-peer/modules/storage/database"
	"github.com/number571/go-peer/settings/testutils"

	payload_adapter "github.com/number571/go-peer/modules/network/anonymity/adapters/payload"
)

const (
	tcPathDBTemplate = "database_test_%d.db"
	tcWait           = 30 * time.Second
	tcIter           = 10
	msgSize          = (100 << 10)
)

func TestComplex(t *testing.T) {
	dbs, nnodes, nodes := testNewNodes(t, tcWait)
	if nodes[0] == nil {
		return
	}
	defer testFreeNodes(dbs[:], nnodes[:], nodes[:])

	wg := sync.WaitGroup{}
	wg.Add(tcIter)

	for i := 0; i < tcIter; i++ {
		go func(i int) {
			defer wg.Done()
			reqBody := fmt.Sprintf("%s (%d)", testutils.TcLargeBody, i)

			// nodes[1] -> nodes[0] -> nodes[2]
			resp, err := nodes[0].Request(
				nodes[1].Queue().Client().PubKey(),
				payload_adapter.NewPayload(testutils.TcHead, []byte(reqBody)),
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

func TestF2FWithoutFriends(t *testing.T) {
	// 5 seconds for wait
	dbs, nnodes, nodes := testNewNodes(t, 2*time.Second)
	if nodes[0] == nil {
		return
	}
	defer testFreeNodes(dbs[:], nnodes[:], nodes[:])

	nodes[0].F2F().Remove(nodes[1].Queue().Client().PubKey())
	nodes[1].F2F().Remove(nodes[0].Queue().Client().PubKey())

	// nodes[1] -> nodes[0] -> nodes[2]
	_, err := nodes[0].Request(
		nodes[1].Queue().Client().PubKey(),
		payload_adapter.NewPayload(testutils.TcHead, []byte(testutils.TcLargeBody)),
	)
	if err != nil {
		return
	}

	t.Errorf("get response without list of friends")
}

// nodes[0], nodes[1] = clients
// nodes[2], nodes[3], nodes[4] = routes
// (nodes[0]) -> nodes[2] -> nodes[3] -> nodes[4] -> (nodes[1])
func testNewNodes(t *testing.T, timeWait time.Duration) ([5]database.IKeyValueDB, [5]network.INode, [5]INode) {
	dbs := [5]database.IKeyValueDB{}
	nnodes := [5]network.INode{}
	nodes := [5]INode{}

	for i := 0; i < 5; i++ {
		dbs[i], nnodes[i], nodes[i] = testNewNode(i, timeWait)
		if nodes[i] == nil {
			t.Errorf("node (%d) is not running", i)
			return [5]database.IKeyValueDB{}, [5]network.INode{}, [5]INode{}
		}
	}

	nodes[0].F2F().Append(nodes[1].Queue().Client().PubKey())
	nodes[1].F2F().Append(nodes[0].Queue().Client().PubKey())

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

	// nodes to routes (nodes[0] -> nodes[2], nodes[1] -> nodes[4])
	nodes[0].Network().Connect(testutils.TgAddrs[2])
	nodes[1].Network().Connect(testutils.TgAddrs[3])

	// routes to routes (nodes[3] -> nodes[2], nodes[3] -> nodes[4])
	nodes[3].Network().Connect(testutils.TgAddrs[2])
	nodes[3].Network().Connect(testutils.TgAddrs[3])

	return dbs, nnodes, nodes
}

func testNewNode(i int, timeWait time.Duration) (database.IKeyValueDB, network.INode, INode) {
	db := database.NewLevelDB(
		database.NewSettings(&database.SSettings{
			FPath:      fmt.Sprintf(tcPathDBTemplate, i),
			FHashing:   true,
			FCipherKey: []byte(testutils.TcKey1),
		}),
	)
	nnode := network.NewNode(
		network.NewSettings(&network.SSettings{
			FCapacity:    (1 << 10),
			FMessageSize: (100 << 10),
			FMaxConnects: 10,
			FTimeWait:    5 * time.Second,
		}),
	)
	node := NewNode(
		NewSettings(&SSettings{
			FRetryEnqueue: 0,
			FTimeWait:     timeWait,
		}),
		db,
		nnode,
		queue.NewQueue(
			queue.NewSettings(&queue.SSettings{
				FCapacity:     10,
				FPullCapacity: 5,
				FDuration:     500 * time.Millisecond,
			}),
			client.NewClient(
				client.NewSettings(&client.SSettings{
					FWorkSize:    10,
					FMessageSize: (100 << 10),
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey),
			),
		),
		friends.NewF2F(),
	)
	if err := node.Run(); err != nil {
		return nil, nil, nil
	}
	return db, nnode, node
}

func testFreeNodes(dbs []database.IKeyValueDB, nnodes []network.INode, nodes []INode) {
	toClose := make([]modules.ICloser, 0, len(dbs)+len(nnodes)+len(nodes))
	for _, db := range dbs {
		toClose = append(toClose, db)
	}
	for _, nnode := range nnodes {
		toClose = append(toClose, nnode)
	}
	for _, node := range nodes {
		toClose = append(toClose, node)
	}
	closer.CloseAll(toClose)
	for i := 0; i < 5; i++ {
		os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, i))
	}
}
