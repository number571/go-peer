package anonymity

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
	testutils "github.com/number571/go-peer/test/_data"

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
			resp, err := nodes[0].FetchPayload(
				nodes[1].GetMessageQueue().GetClient().GetPubKey(),
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

	nodes[0].GetListPubKeys().DelPubKey(nodes[1].GetMessageQueue().GetClient().GetPubKey())
	nodes[1].GetListPubKeys().DelPubKey(nodes[0].GetMessageQueue().GetClient().GetPubKey())

	// nodes[1] -> nodes[0] -> nodes[2]
	_, err := nodes[0].FetchPayload(
		nodes[1].GetMessageQueue().GetClient().GetPubKey(),
		NewPayload(testutils.TcHead, []byte(testutils.TcLargeBody)),
	)
	if err != nil {
		return
	}

	t.Error("get response without list of friends")
}

// nodes[0], nodes[1] = clients
// nodes[2], nodes[3], nodes[4] = routes
// (nodes[0]) -> nodes[2] -> nodes[3] -> nodes[4] -> (nodes[1])
func testNewNodes(t *testing.T, timeWait time.Duration) [5]INode {
	nodes := [5]INode{}
	addrs := [5]string{"", "", testutils.TgAddrs[2], "", testutils.TgAddrs[3]}

	for i := 0; i < 5; i++ {
		nodes[i] = testNewNode(i, timeWait, addrs[i])
		if nodes[i] == nil {
			t.Errorf("node (%d) is not running", i)
			return [5]INode{}
		}
	}

	nodes[0].GetListPubKeys().AddPubKey(nodes[1].GetMessageQueue().GetClient().GetPubKey())
	nodes[1].GetListPubKeys().AddPubKey(nodes[0].GetMessageQueue().GetClient().GetPubKey())

	for _, node := range nodes {
		node.HandleFunc(
			testutils.TcHead,
			func(_ INode, _ asymmetric.IPubKey, _, reqBytes []byte) []byte {
				// send response
				return []byte(fmt.Sprintf("%s (response)", string(reqBytes)))
			},
		)
	}

	if err := nodes[2].GetNetworkNode().Run(); err != nil {
		panic(err)
	}
	if err := nodes[4].GetNetworkNode().Run(); err != nil {
		panic(err)
	}

	time.Sleep(200 * time.Millisecond)

	// nodes to routes (nodes[0] -> nodes[2], nodes[1] -> nodes[4])
	nodes[0].GetNetworkNode().AddConnect(testutils.TgAddrs[2])
	nodes[1].GetNetworkNode().AddConnect(testutils.TgAddrs[3])

	// routes to routes (nodes[3] -> nodes[2], nodes[3] -> nodes[4])
	nodes[3].GetNetworkNode().AddConnect(testutils.TgAddrs[2])
	nodes[3].GetNetworkNode().AddConnect(testutils.TgAddrs[3])

	return nodes
}

func testNewNode(i int, timeWait time.Duration, addr string) INode {
	node := NewNode(
		NewSettings(&SSettings{
			FRetryEnqueue: 0,
			FTimeWait:     timeWait,
		}),
		logger.NewLogger(logger.NewSettings(&logger.SSettings{})),
		NewWrapperDB().Set(database.NewLevelDB(
			database.NewSettings(&database.SSettings{
				FPath:      fmt.Sprintf(tcPathDBTemplate, i),
				FHashing:   true,
				FCipherKey: []byte(testutils.TcKey1),
			}),
		)),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:     addr,
				FCapacity:    (1 << 10),
				FMaxConnects: 10,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FMessageSize: (100 << 10),
					FTimeWait:    5 * time.Second,
				}),
			}),
		),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{
				FCapacity:     10,
				FPullCapacity: 5,
				FDuration:     500 * time.Millisecond,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FWorkSize:    10,
					FMessageSize: (100 << 10),
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey),
			),
		),
		asymmetric.NewListPubKeys(),
	)
	if err := node.Run(); err != nil {
		return nil
	}
	return node
}

func testFreeNodes(nodes []INode) {
	for _, node := range nodes {
		node.GetWrapperDB().Close()
		types.StopAll([]types.ICommand{node, node.GetNetworkNode()})
	}
	for i := 0; i < 5; i++ {
		os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, i))
	}
}
