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

	"github.com/number571/go-peer/pkg/network/anonymity/adapters"
	"github.com/number571/go-peer/pkg/network/conn"
)

const (
	tcPathDBTemplate = "database_test_%d_%d.db"
	tcWait           = time.Minute
	tcIter           = 10
)

func TestComplex(t *testing.T) {
	for i := 0; i < 5; i++ {
		os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, i, 0))
	}

	nodes := testNewNodes(t, tcWait, 0)
	if nodes[0] == nil {
		return
	}
	defer testFreeNodes(nodes[:], 0)

	wg := sync.WaitGroup{}
	wg.Add(tcIter)

	for i := 0; i < tcIter; i++ {
		go func(i int) {
			defer wg.Done()
			reqBody := fmt.Sprintf("%s (%d)", testutils.TcLargeBody, i)

			// nodes[1] -> nodes[0] -> nodes[2]
			resp, err := nodes[0].FetchPayload(
				nodes[1].GetMessageQueue().GetClient().GetPubKey(),
				adapters.NewPayload(testutils.TcHead, []byte(reqBody)),
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
	for i := 0; i < 5; i++ {
		os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, i, 1))
	}

	// 3 seconds for wait
	nodes := testNewNodes(t, 3*time.Second, 1)
	if nodes[0] == nil {
		t.Error("[f2f] can't create node")
		return
	}
	defer testFreeNodes(nodes[:], 1)

	nodes[0].GetListPubKeys().DelPubKey(nodes[1].GetMessageQueue().GetClient().GetPubKey())
	nodes[1].GetListPubKeys().DelPubKey(nodes[0].GetMessageQueue().GetClient().GetPubKey())

	// nodes[1] -> nodes[0] -> nodes[2]
	_, err := nodes[0].FetchPayload(
		nodes[1].GetMessageQueue().GetClient().GetPubKey(),
		adapters.NewPayload(testutils.TcHead, []byte(testutils.TcLargeBody)),
	)
	if err != nil {
		return
	}

	t.Error("get response without list of friends")
}

// nodes[0], nodes[1] = clients
// nodes[2], nodes[3], nodes[4] = routes
// (nodes[0]) -> nodes[2] -> nodes[3] -> nodes[4] -> (nodes[1])
func testNewNodes(t *testing.T, timeWait time.Duration, typeDB int) [5]INode {
	nodes := [5]INode{}
	addrs := [5]string{"", "", testutils.TgAddrs[2], "", testutils.TgAddrs[3]}

	for i := 0; i < 5; i++ {
		nodes[i] = testNewNode(i, timeWait, addrs[i], typeDB)
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
		t.Error(err)
		return [5]INode{}
	}
	if err := nodes[4].GetNetworkNode().Run(); err != nil {
		t.Error(err)
		return [5]INode{}
	}

	time.Sleep(time.Second)

	// nodes to routes (nodes[0] -> nodes[2], nodes[1] -> nodes[4])
	nodes[0].GetNetworkNode().AddConnect(testutils.TgAddrs[2])
	nodes[1].GetNetworkNode().AddConnect(testutils.TgAddrs[3])

	// routes to routes (nodes[3] -> nodes[2], nodes[3] -> nodes[4])
	nodes[3].GetNetworkNode().AddConnect(testutils.TgAddrs[2])
	nodes[3].GetNetworkNode().AddConnect(testutils.TgAddrs[3])

	return nodes
}

func testNewNode(i int, timeWait time.Duration, addr string, typeDB int) INode {
	db, err := database.NewSQLiteDB(
		database.NewSettings(&database.SSettings{
			FPath:      fmt.Sprintf(tcPathDBTemplate, i, typeDB),
			FHashing:   true,
			FCipherKey: []byte("CIPHER"),
		}),
	)
	if err != nil {
		return nil
	}
	node := NewNode(
		NewSettings(&SSettings{
			FServiceName:   "TEST",
			FRetryEnqueue:  0,
			FNetworkMask:   1,
			FFetchTimeWait: timeWait,
		}),
		logger.NewLogger(logger.NewSettings(&logger.SSettings{})),
		NewWrapperDB().Set(db),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:     addr,
				FCapacity:    testutils.TCCapacity,
				FMaxConnects: testutils.TCMaxConnects,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FMessageSize:   testutils.TCMessageSize,
					FLimitVoidSize: 1, // not used
					FFetchTimeWait: 1, // not used
				}),
			}),
		),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{
				FMainCapacity: testutils.TCQueueCapacity,
				FPoolCapacity: testutils.TCQueueCapacity,
				FDuration:     time.Second,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FWorkSize:    testutils.TCWorkSize,
					FMessageSize: testutils.TCMessageSize,
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

func testFreeNodes(nodes []INode, typeDB int) {
	for _, node := range nodes {
		node.GetWrapperDB().Close()
		types.StopAll([]types.ICommand{node, node.GetNetworkNode()})
	}
	for i := 0; i < 5; i++ {
		os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, i, typeDB))
	}
}
