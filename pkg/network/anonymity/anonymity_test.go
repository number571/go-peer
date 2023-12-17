package anonymity

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/queue_set"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/storage/database"
	testutils "github.com/number571/go-peer/test/_data"

	"github.com/number571/go-peer/pkg/network/anonymity/adapters"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

const (
	tcPathDBTemplate = "database_test_%d_%d.db"
	tcIter           = 10
)

func TestNodeSettings(t *testing.T) {
	t.Parallel()

	node, cancels := testNewNode(time.Minute, "", 9, 0)
	defer testFreeNodes([]INode{node}, []context.CancelFunc{cancels}, 9)

	sett := node.GetSettings()
	if sett.GetFetchTimeWait() != time.Minute {
		t.Error("sett.GetFetchTimeWait() != time.Minute")
		return
	}
	_ = node.GetLogger()
}

func TestSettings(t *testing.T) {
	t.Parallel()

	for i := 0; i < 3; i++ {
		testSettings(t, i)
	}
}

func testSettings(t *testing.T, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	switch n {
	case 0:
		_ = NewSettings(&SSettings{
			FRetryEnqueue:  0,
			FNetworkMask:   1,
			FFetchTimeWait: time.Second,
		})
	case 1:
		_ = NewSettings(&SSettings{
			FServiceName:   "TEST",
			FRetryEnqueue:  0,
			FFetchTimeWait: time.Second,
		})
	case 2:
		_ = NewSettings(&SSettings{
			FServiceName:  "TEST",
			FRetryEnqueue: 0,
			FNetworkMask:  1,
		})
	}
}

func TestComplexFetchPayload(t *testing.T) {
	t.Parallel()

	addresses := [2]string{testutils.TgAddrs[2], testutils.TgAddrs[3]}
	nodes, cancels := testNewNodes(t, time.Minute, addresses, 0)
	if nodes[0] == nil {
		t.Error("nodes is null")
		return
	}
	defer testFreeNodes(nodes[:], cancels[:], 0)

	wg := sync.WaitGroup{}
	wg.Add(tcIter)

	for i := 0; i < tcIter; i++ {
		go func(i int) {
			defer wg.Done()
			reqBody := fmt.Sprintf("%s (%d)", testutils.TcBody, i)

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
	t.Parallel()

	// 3 seconds for wait
	addresses := [2]string{testutils.TgAddrs[31], testutils.TgAddrs[32]}
	nodes, cancels := testNewNodes(t, 3*time.Second, addresses, 1)
	if nodes[0] == nil {
		t.Error("nodes is null")
		return
	}
	defer testFreeNodes(nodes[:], cancels[:], 1)

	nodes[0].GetListPubKeys().DelPubKey(nodes[1].GetMessageQueue().GetClient().GetPubKey())
	nodes[1].GetListPubKeys().DelPubKey(nodes[0].GetMessageQueue().GetClient().GetPubKey())

	// nodes[1] -> nodes[0] -> nodes[2]
	_, err := nodes[0].FetchPayload(
		nodes[1].GetMessageQueue().GetClient().GetPubKey(),
		adapters.NewPayload(testutils.TcHead, []byte(testutils.TcBody)),
	)
	if err != nil {
		return
	}

	t.Error("get response without list of friends")
}

func TestDataType(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	switch {
	case isRequest([]byte{}):
		t.Error("is request = true with void bytes")
	case isResponse([]byte{}):
		t.Error("is response = true with void bytes")
	}

	_ = unwrapBytes([]byte{})
}

func TestWrapper(t *testing.T) {
	t.Parallel()

	wrapper := NewWrapperDB()
	if db := wrapper.Get(); db != nil {
		t.Error("db is not null")
		return
	}
	if err := wrapper.Close(); err != nil {
		t.Error(err)
		return
	}
}

func TestFetchPayload(t *testing.T) {
	t.Parallel()

	addresses := [2]string{testutils.TgAddrs[35], testutils.TgAddrs[36]}
	nodes, cancels := testNewNodes(t, time.Minute, addresses, 4)
	if nodes[0] == nil {
		t.Error("nodes is null")
		return
	}
	defer testFreeNodes(nodes[:], cancels[:], 4)

	nodes[1].HandleFunc(
		testutils.TcHead,
		func(_ INode, _ asymmetric.IPubKey, reqBytes []byte) ([]byte, error) {
			return []byte(fmt.Sprintf("echo: '%s'", string(reqBytes))), nil
		},
	)

	msgBody := "hello, world!"
	result, err := nodes[0].FetchPayload(
		nodes[1].GetMessageQueue().GetClient().GetPubKey(),
		adapters.NewPayload(testutils.TcHead, []byte(msgBody)),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if string(result) != fmt.Sprintf("echo: '%s'", msgBody) {
		t.Error("got invalid message body")
		return
	}

	if err := nodes[0].GetNetworkNode().DelConnection(testutils.TgAddrs[35]); err != nil {
		t.Error(err)
		return
	}

	_, err2 := nodes[0].FetchPayload(
		nodes[1].GetMessageQueue().GetClient().GetPubKey(),
		adapters.NewPayload(testutils.TcHead, []byte(msgBody)),
	)
	if err2 == nil {
		t.Error("success fetch payload without connections")
		return
	}
}

func TestBroadcastPayload(t *testing.T) {
	t.Parallel()

	addresses := [2]string{testutils.TgAddrs[33], testutils.TgAddrs[34]}
	nodes, cancels := testNewNodes(t, time.Minute, addresses, 3)
	if nodes[0] == nil {
		t.Error("nodes is null")
		return
	}
	defer testFreeNodes(nodes[:], cancels[:], 3)

	chResult := make(chan string)
	nodes[1].HandleFunc(
		testutils.TcHead,
		func(_ INode, _ asymmetric.IPubKey, reqBytes []byte) ([]byte, error) {
			res := fmt.Sprintf("echo: '%s'", string(reqBytes))
			go func() { chResult <- res }()
			return nil, nil
		},
	)

	msgBody := "hello, world!"
	err := nodes[0].BroadcastPayload(
		nodes[1].GetMessageQueue().GetClient().GetPubKey(),
		adapters.NewPayload(testutils.TcHead, []byte(msgBody)),
	)
	if err != nil {
		t.Error(err)
		return
	}

	select {
	case x := <-chResult:
		if x != fmt.Sprintf("echo: '%s'", msgBody) {
			t.Error("got invalid message body")
			return
		}
		// success
	case <-time.After(5 * time.Second):
		t.Error("error: time after 5 seconds")
		return
	}

	if err := nodes[0].GetNetworkNode().DelConnection(testutils.TgAddrs[33]); err != nil {
		t.Error(err)
		return
	}

	err2 := nodes[0].BroadcastPayload(
		nodes[1].GetMessageQueue().GetClient().GetPubKey(),
		adapters.NewPayload(testutils.TcHead, []byte(msgBody)),
	)
	if err2 == nil {
		t.Error("success broadcast payload without connections")
		return
	}
}

// func TestEnqueuePayload(t *testing.T) {
// 	t.Parallel()

// 	addresses := [2]string{testutils.TgAddrs[38], testutils.TgAddrs[39]}
// 	nodes := testNewNodes(t, time.Minute, addresses, 8)
// 	if nodes[0] == nil {
// 		t.Error("nodes is null")
// 		return
// 	}
// 	defer testFreeNodes(nodes[:], 8)

// 	node := nodes[0].(*sNode)
// 	client := nodes[0].GetMessageQueue().GetClient()
// 	pubKey := nodes[1].GetMessageQueue().GetClient().GetPubKey()

// 	msgBody := "hello, world!"
// 	pld := adapters.NewPayload(testutils.TcHead, []byte(msgBody)).ToOrigin()
// 	if err := node.enqueuePayload('?', pubKey, pld); err == nil {
// 		t.Error("success with undefined type of message")
// 		return
// 	}

// 	overheadBody := random.NewStdPRNG().GetBytes(testutils.TCMessageSize + 1)
// 	overPld := adapters.NewPayload(testutils.TcHead, overheadBody).ToOrigin()
// 	if err := node.enqueuePayload(cIsRequest, pubKey, overPld); err == nil {
// 		t.Error("success with overhead message")
// 		return
// 	}

// 	msg, err := client.EncryptPayload(
// 		pubKey,
// 		adapters.NewPayload(
// 			testutils.TcHead,
// 			wrapRequest([]byte(msgBody)),
// 		).ToOrigin(),
// 	)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	for i := 0; i < testutils.TCQueueCapacity; i++ {
// 		if err := node.send(msg); err != nil {
// 			t.Error("failed send message (push to queue)")
// 			return
// 		}
// 	}

// 	for i := 0; i < 5; i++ {
// 		if err := node.enqueuePayload(cIsRequest, pubKey, pld); err != nil {
// 			return
// 		}
// 		for j := 0; j < testutils.TCQueueCapacity; j++ {
// 			if err := node.send(msg); err != nil {
// 				break // fill queue on limit
// 			}
// 		}
// 	}

// 	t.Error("success enqueue payload over queue capacity")
// }

func TestHandleWrapper(t *testing.T) {
	t.Parallel()

	_node, cancel := testNewNode(time.Minute, "", 7, 0)
	defer testFreeNodes([]INode{_node}, []context.CancelFunc{cancel}, 7)

	node := _node.(*sNode)
	handler := node.handleWrapper()
	client := node.fQueue.GetClient()
	pubKey := client.GetPubKey()

	node.GetListPubKeys().AddPubKey(pubKey)

	msgBody := "hello, world!"
	msg, err := client.EncryptPayload(
		pubKey,
		adapters.NewPayload(
			testutils.TcHead,
			wrapRequest([]byte(msgBody)),
		).ToOrigin(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	sett := net_message.NewSettings(&net_message.SSettings{})
	netMsg := node.testNewNetworkMessage(sett, msg)
	if err := handler(nil, nil, netMsg); err != nil {
		t.Error(err)
		return
	}

	if err := handler(nil, nil, netMsg); err != nil {
		t.Error("repeated message:", err.Error())
		return
	}

	node.HandleFunc(
		111,
		func(_ INode, _ asymmetric.IPubKey, _ []byte) ([]byte, error) {
			return nil, errors.New("some error")
		},
	)

	msg2, err := client.EncryptPayload(
		pubKey,
		adapters.NewPayload(
			111,
			wrapRequest([]byte(msgBody)),
		).ToOrigin(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	netMsg2 := node.testNewNetworkMessage(sett, msg2)
	if err := handler(nil, nil, netMsg2); err != nil {
		t.Error(err) // works only logger
		return
	}

	msg3, err := client.EncryptPayload(
		pubKey,
		adapters.NewPayload(
			111,
			[]byte("?"+msgBody),
		).ToOrigin(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	netMsg3 := node.testNewNetworkMessage(sett, msg3)
	if err := handler(nil, nil, netMsg3); err != nil {
		t.Error(err) // works only logger
		return
	}

	msg4, err := client.EncryptPayload(
		pubKey,
		adapters.NewPayload(
			111,
			wrapResponse([]byte(msgBody)),
		).ToOrigin(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	netMsg4 := node.testNewNetworkMessage(sett, msg4)
	if err := handler(nil, nil, netMsg4); err != nil {
		t.Error(err) // works only logger
		return
	}
}

func TestStoreHashWithBroadcastMessage(t *testing.T) {
	t.Parallel()

	_node, cancel := testNewNode(time.Minute, "", 6, 0)
	defer testFreeNodes([]INode{_node}, []context.CancelFunc{cancel}, 6)

	node := _node.(*sNode)
	client := node.fQueue.GetClient()

	msgBody := "hello, world!"
	msg, err := client.EncryptPayload(
		client.GetPubKey(),
		adapters.NewPayload(
			testutils.TcHead,
			wrapRequest([]byte(msgBody)),
		).ToOrigin(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	sett := net_message.NewSettings(&net_message.SSettings{})
	netMsg := node.testNewNetworkMessage(sett, msg)
	logBuilder := anon_logger.NewLogBuilder("_")

	if ok, err := node.storeHashWithBroadcast(logBuilder, netMsg); !ok || err != nil {
		t.Error(err)
		return
	}

	if ok, err := node.storeHashWithBroadcast(logBuilder, netMsg); ok || err != nil {
		switch {
		case ok:
			t.Error("success store one message again")
		case err != nil:
			t.Error("got error with try store twice same message")
		}
		return
	}

	node.GetWrapperDB().Set(nil)
	if ok, err := node.storeHashWithBroadcast(logBuilder, netMsg); ok || err == nil {
		t.Error("success use store function with null database")
		return
	}

	time.Sleep(time.Second + 100*time.Millisecond) // queue duration
}

func TestRecvSendMessage(t *testing.T) {
	t.Parallel()

	_node, cancel := testNewNode(time.Minute, "", 5, 0)
	defer testFreeNodes([]INode{_node}, []context.CancelFunc{cancel}, 5)

	node := _node.(*sNode)
	if _, err := node.recv("not_exist"); err == nil {
		t.Error("success got action by undefined key")
		return
	}

	client := node.fQueue.GetClient()
	pubKey := client.GetPubKey()
	actionKey := newActionKey(pubKey, 111)

	node.setAction(actionKey)
	action, ok := node.getAction(actionKey)
	if !ok {
		t.Error("undefined created action key")
		return
	}

	close(action)

	if _, err := node.recv(actionKey); err == nil {
		t.Error("success got closed action")
		return
	}

	msgBody := "hello, world!"
	msg, err := client.EncryptPayload(
		pubKey,
		adapters.NewPayload(
			testutils.TcHead,
			wrapRequest([]byte(msgBody)),
		).ToOrigin(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < testutils.TCQueueCapacity; i++ {
		if err := node.send(msg); err != nil {
			t.Error("failed send message (push to queue)")
			return
		}
	}

	for i := 0; i < 3; i++ {
		// message can be dequeued in the send's call time
		if err := node.send(msg); err != nil {
			return
		}
	}

	t.Error("success send message (push to queue) over queue capacity")
}

// nodes[0], nodes[1] = clients
// nodes[2], nodes[3], nodes[4] = routes
// nodes[2], nodes[4] are have open ports
// Scheme: (nodes[0]) -> nodes[2] -> nodes[3] -> nodes[4] -> (nodes[1])
func testNewNodes(t *testing.T, timeWait time.Duration, addresses [2]string, typeDB int) ([5]INode, [5]context.CancelFunc) {
	nodes := [5]INode{}
	cancels := [5]context.CancelFunc{}
	addrs := [5]string{"", "", addresses[0], "", addresses[1]}

	for i := 0; i < 5; i++ {
		nodes[i], cancels[i] = testNewNode(timeWait, addrs[i], typeDB, i)
		if nodes[i] == nil {
			t.Errorf("node (%d) is not running %d", i, typeDB)
			return [5]INode{}, [5]context.CancelFunc{}
		}
	}

	nodes[0].GetListPubKeys().AddPubKey(nodes[1].GetMessageQueue().GetClient().GetPubKey())
	nodes[1].GetListPubKeys().AddPubKey(nodes[0].GetMessageQueue().GetClient().GetPubKey())

	for _, node := range nodes {
		node.HandleFunc(
			testutils.TcHead,
			func(_ INode, _ asymmetric.IPubKey, reqBytes []byte) ([]byte, error) {
				// send response
				return []byte(fmt.Sprintf("%s (response)", string(reqBytes))), nil
			},
		)
	}

	if err := nodes[2].GetNetworkNode().Listen(); err != nil {
		t.Error(err)
		return [5]INode{}, [5]context.CancelFunc{}
	}
	if err := nodes[4].GetNetworkNode().Listen(); err != nil {
		t.Error(err)
		return [5]INode{}, [5]context.CancelFunc{}
	}

	time.Sleep(200 * time.Millisecond)

	// nodes to routes (nodes[0] -> nodes[2], nodes[1] -> nodes[4])
	if err := nodes[0].GetNetworkNode().AddConnection(addresses[0]); err != nil {
		t.Error(err)
		return [5]INode{}, [5]context.CancelFunc{}
	}
	if err := nodes[1].GetNetworkNode().AddConnection(addresses[1]); err != nil {
		t.Error(err)
		return [5]INode{}, [5]context.CancelFunc{}
	}

	// routes to routes (nodes[3] -> nodes[2], nodes[3] -> nodes[4])
	if err := nodes[3].GetNetworkNode().AddConnection(addresses[0]); err != nil {
		t.Error(err)
		return [5]INode{}, [5]context.CancelFunc{}
	}
	if err := nodes[3].GetNetworkNode().AddConnection(addresses[1]); err != nil {
		t.Error(err)
		return [5]INode{}, [5]context.CancelFunc{}
	}

	return nodes, cancels
}

func testNewNode(timeWait time.Duration, addr string, typeDB, numDB int) (INode, context.CancelFunc) {
	db, err := database.NewKVDatabase(
		storage.NewSettings(&storage.SSettings{
			FPath:     fmt.Sprintf(tcPathDBTemplate, typeDB, numDB),
			FWorkSize: testutils.TCWorkSize,
			FPassword: "CIPHER",
		}),
	)
	if err != nil {
		return nil, nil
	}
	networkMask := uint64(1)
	networkKey := "old_network_key"
	node := NewNode(
		NewSettings(&SSettings{
			FServiceName:   "TEST",
			FNetworkMask:   networkMask,
			FFetchTimeWait: timeWait,
		}),
		logger.NewLogger(
			logger.NewSettings(&logger.SSettings{}),
			func(_ logger.ILogArg) string { return "" },
		),
		NewWrapperDB().Set(db),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:      addr,
				FMaxConnects:  testutils.TCMaxConnects,
				FReadTimeout:  timeWait,
				FWriteTimeout: timeWait,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FNetworkKey:       networkKey,
					FWorkSizeBits:     testutils.TCWorkSize,
					FMessageSizeBytes: testutils.TCMessageSize,
					FWaitReadDeadline: time.Hour,
					FReadDeadline:     time.Minute,
					FWriteDeadline:    time.Minute,
				}),
			}),
			queue_set.NewQueueSet(
				queue_set.NewSettings(&queue_set.SSettings{
					FCapacity: testutils.TCCapacity,
				}),
			),
		),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{
				FMainCapacity: testutils.TCQueueCapacity,
				FPoolCapacity: testutils.TCQueueCapacity,
				FDuration:     time.Second,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FMessageSizeBytes: testutils.TCMessageSize,
					FKeySizeBits:      testutils.TcKeySize,
				}),
				asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
			),
		).WithNetworkSettings(
			networkMask,
			net_message.NewSettings(&net_message.SSettings{
				FNetworkKey:   networkKey,
				FWorkSizeBits: testutils.TCWorkSize,
			}),
		),
		asymmetric.NewListPubKeys(),
	)
	ctx, cancel := context.WithCancel(context.Background())
	go func() { _ = node.Run(ctx) }()
	return node, cancel
}

func testFreeNodes(nodes []INode, cancels []context.CancelFunc, typeDB int) {
	for i, node := range nodes {
		node.GetWrapperDB().Close()
		node.GetNetworkNode().Close()
		cancels[i]()
	}
	testDeleteDB(typeDB)
}

func testDeleteDB(typeDB int) {
	for i := 0; i < 5; i++ {
		os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, typeDB, i))
	}
}

func (p *sNode) testNewNetworkMessage(pSett net_message.ISettings, pMsg message.IMessage) net_message.IMessage {
	return net_message.NewMessage(
		pSett,
		payload.NewPayload(
			p.fSettings.GetNetworkMask(),
			pMsg.ToBytes(),
		),
	)
}
