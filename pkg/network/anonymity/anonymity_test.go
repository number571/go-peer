package anonymity

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/number571/go-peer/internal/testutils"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/closer"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/friends"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/queue"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"

	"github.com/number571/go-peer/pkg/network/conn"
)

const (
	tcPathDBTemplate = "database_test_%d.db"
	tcWait           = 30 * time.Second
	tcIter           = 10
	msgSize          = (100 << 10)
)

func TestComplex(t *testing.T) {
	nodes := testNewNodes(t, tcWait)
	if nodes[0] == nil {
		return
	}
	defer testFreeNodes(nodes[:])

	wg := sync.WaitGroup{}
	wg.Add(tcIter)

	for i := 0; i < tcIter; i++ {
		go func(i int) {
			defer wg.Done()
			reqBody := fmt.Sprintf("%s (%d)", testutils.TcLargeBody, i)

			// nodes[1] -> nodes[0] -> nodes[2]
			resp, err := nodes[0].Request(
				nodes[1].Queue().Client().PubKey(),
				NewPayload(testutils.TcHead, []byte(reqBody)),
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
	// 3 seconds for wait
	nodes := testNewNodes(t, 3*time.Second)
	if nodes[0] == nil {
		return
	}
	defer testFreeNodes(nodes[:])

	nodes[0].F2F().Remove(nodes[1].Queue().Client().PubKey())
	nodes[1].F2F().Remove(nodes[0].Queue().Client().PubKey())

	// nodes[1] -> nodes[0] -> nodes[2]
	_, err := nodes[0].Request(
		nodes[1].Queue().Client().PubKey(),
		NewPayload(testutils.TcHead, []byte(testutils.TcLargeBody)),
	)
	if err != nil {
		return
	}

	t.Errorf("get response without list of friends")
}

// nodes[0], nodes[1] = clients
// nodes[2], nodes[3], nodes[4] = routes
// (nodes[0]) -> nodes[2] -> nodes[3] -> nodes[4] -> (nodes[1])
func testNewNodes(t *testing.T, timeWait time.Duration) [5]INode {
	nodes := [5]INode{}

	for i := 0; i < 5; i++ {
		nodes[i] = testNewNode(i, timeWait)
		if nodes[i] == nil {
			t.Errorf("node (%d) is not running", i)
			return [5]INode{}
		}
	}

	nodes[0].F2F().Append(nodes[1].Queue().Client().PubKey())
	nodes[1].F2F().Append(nodes[0].Queue().Client().PubKey())

	for _, node := range nodes {
		node.Handle(
			testutils.TcHead,
			func(node INode, sender asymmetric.IPubKey, reqBytes []byte) []byte {
				// send response
				return []byte(fmt.Sprintf("%s (response)", string(reqBytes)))
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

	return nodes
}

func testNewNode(i int, timeWait time.Duration) INode {
	node := NewNode(
		NewSettings(&SSettings{
			FRetryEnqueue: 0,
			FTimeWait:     timeWait,
		}),
		logger.NewLogger(logger.NewSettings(&logger.SSettings{})),
		database.NewLevelDB(
			database.NewSettings(&database.SSettings{
				FHashing:   true,
				FCipherKey: []byte(testutils.TcKey1),
			}),
			fmt.Sprintf(tcPathDBTemplate, i),
		),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FCapacity:    (1 << 10),
				FMaxConnects: 10,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FMessageSize: (100 << 10),
					FTimeWait:    5 * time.Second,
				}),
			}),
		),
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
		return nil
	}
	return node
}

func testFreeNodes(nodes []INode) {
	toClose := make([]types.ICloser, 0, len(nodes)*3)
	for _, node := range nodes {
		toClose = append(toClose, node, node.KeyValueDB(), node.Network())
	}
	closer.CloseAll(toClose)
	for i := 0; i < 5; i++ {
		os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, i))
	}
}
