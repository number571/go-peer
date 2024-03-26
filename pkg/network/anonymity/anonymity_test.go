package anonymity

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/cache/lru"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/database"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/utils"

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

	node, cancels := testNewNode(time.Minute, "", 9, 0, 0, false)
	defer testFreeNodes([]INode{node}, []context.CancelFunc{cancels}, 9)

	sett := node.GetSettings()
	if sett.GetFetchTimeout() != time.Minute {
		t.Error("sett.GetFetchTimeout() != time.Minute")
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
			FRetryEnqueue: 0,
			FNetworkMask:  1,
			FFetchTimeout: time.Second,
		})
	case 1:
		_ = NewSettings(&SSettings{
			FServiceName:  "TEST",
			FRetryEnqueue: 0,
			FFetchTimeout: time.Second,
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

	ctx := context.Background()

	for i := 0; i < tcIter; i++ {
		go func(i int) {
			defer wg.Done()
			reqBody := fmt.Sprintf("%s (%d)", testutils.TcBody, i)

			// nodes[1] -> nodes[0] -> nodes[2]
			resp, err := nodes[0].FetchPayload(
				ctx,
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

	ctx := context.Background()

	// nodes[1] -> nodes[0] -> nodes[2]
	_, err := nodes[0].FetchPayload(
		ctx,
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

	wrapper := NewDBWrapper()
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
		func(_ context.Context, _ INode, _ asymmetric.IPubKey, reqBytes []byte) ([]byte, error) {
			return []byte(fmt.Sprintf("echo: '%s'", string(reqBytes))), nil
		},
	)

	ctx := context.Background()
	_, err := nodes[0].FetchPayload(
		ctx,
		nodes[1].GetMessageQueue().GetClient().GetPubKey(),
		adapters.NewPayload(testutils.TcHead, []byte(testutils.TcLargeBody)),
	)
	if err == nil {
		t.Error("success fetch payload with large body")
		return
	}

	msgBody := "hello, world!"
	result, err1 := nodes[0].FetchPayload(
		ctx,
		nodes[1].GetMessageQueue().GetClient().GetPubKey(),
		adapters.NewPayload(testutils.TcHead, []byte(msgBody)),
	)
	if err1 != nil {
		t.Error(err1)
		return
	}

	if string(result) != fmt.Sprintf("echo: '%s'", msgBody) {
		t.Error("got invalid message body")
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
		func(_ context.Context, _ INode, _ asymmetric.IPubKey, reqBytes []byte) ([]byte, error) {
			res := fmt.Sprintf("echo: '%s'", string(reqBytes))
			go func() { chResult <- res }()
			return nil, nil
		},
	)

	ctx := context.Background()
	err := nodes[0].SendPayload(
		ctx,
		nodes[1].GetMessageQueue().GetClient().GetPubKey(),
		payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcLargeBody)),
	)
	if err == nil {
		t.Error("success broadcast payload with large body")
		return
	}

	msgBody := "hello, world!"
	err1 := nodes[0].SendPayload(
		ctx,
		nodes[1].GetMessageQueue().GetClient().GetPubKey(),
		payload.NewPayload(uint64(testutils.TcHead), []byte(msgBody)),
	)
	if err1 != nil {
		t.Error(err1)
		return
	}

	select {
	case x := <-chResult:
		if x != fmt.Sprintf("echo: '%s'", msgBody) {
			t.Error("got invalid message body")
			return
		}
		// success
	case <-time.After(time.Minute):
		t.Error("error: time after 1 minute")
		return
	}
}

func TestPanicEnqueuePayload(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	node, cancels := testNewNode(time.Minute, "", 12, 0, 0, false)
	defer testFreeNodes([]INode{node}, []context.CancelFunc{cancels}, 12)

	sNode := node.(*sNode)

	ctx := context.Background()
	logBuilder := anon_logger.NewLogBuilder("test")

	msgBody := "hello, world!"
	pld := adapters.NewPayload(testutils.TcHead, []byte(msgBody)).ToOrigin()

	pubKey := sNode.GetMessageQueue().GetClient().GetPubKey()
	_ = sNode.enqueuePayload(ctx, logBuilder, '?', pubKey, pld)
}

func TestEnqueuePayload(t *testing.T) {
	t.Parallel()

	addresses := [2]string{testutils.TgAddrs[38], testutils.TgAddrs[39]}
	nodes, cancels := testNewNodes(t, time.Minute, addresses, 8)
	if nodes[0] == nil {
		t.Error("nodes is null")
		return
	}
	defer testFreeNodes(nodes[:], cancels[:], 8)

	node := nodes[0].(*sNode)
	client := nodes[0].GetMessageQueue().GetClient()
	pubKey := nodes[1].GetMessageQueue().GetClient().GetPubKey()

	ctx := context.Background()
	logBuilder := anon_logger.NewLogBuilder("test")

	msgBody := "hello, world!"
	pld := adapters.NewPayload(testutils.TcHead, []byte(msgBody)).ToOrigin()

	overheadBody := random.NewStdPRNG().GetBytes(testutils.TCMessageSize + 1)
	overPld := adapters.NewPayload(testutils.TcHead, overheadBody).ToOrigin()
	if err := node.enqueuePayload(ctx, logBuilder, cIsRequest, pubKey, overPld); err == nil {
		t.Error("success with overhead message")
		return
	}

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
		if err := node.enqueueMessage(ctx, msg); err != nil {
			t.Error("failed send message (push to queue)")
			return
		}
	}

	// after full queue
	for i := 0; i < 2*testutils.TCQueueCapacity; i++ {
		if err := node.enqueuePayload(ctx, logBuilder, cIsRequest, pubKey, pld); err != nil {
			return
		}
	}

	t.Error("success enqueue payload over queue capacity")
}

func TestHandleWrapper(t *testing.T) {
	t.Parallel()

	_node, cancel := testNewNode(time.Minute, "", 7, 0, 0, true)
	defer testFreeNodes([]INode{_node}, []context.CancelFunc{cancel}, 7)

	node := _node.(*sNode)
	handler := node.handleWrapper()
	client := node.fQueue.GetClient()
	pubKey := client.GetPubKey()

	// // ignore add public key (f2f_disabled=true)
	// node.GetListPubKeys().AddPubKey(pubKey)

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

	ctx := context.Background()

	sett := net_message.NewSettings(&net_message.SSettings{})
	netMsg := node.testNewNetworkMessage(sett, msg.ToBytes())
	if err := handler(ctx, nil, nil, netMsg); err != nil {
		t.Error(err)
		return
	}

	if err := handler(ctx, nil, nil, netMsg); err != nil {
		t.Error("repeated message:", err.Error())
		return
	}

	node.HandleFunc(
		111,
		func(_ context.Context, _ INode, _ asymmetric.IPubKey, _ []byte) ([]byte, error) {
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

	netMsg2 := node.testNewNetworkMessage(sett, msg2.ToBytes())
	if err := handler(ctx, nil, nil, netMsg2); err != nil {
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

	netMsg3 := node.testNewNetworkMessage(sett, msg3.ToBytes())
	if err := handler(ctx, nil, nil, netMsg3); err != nil {
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

	netMsg4 := node.testNewNetworkMessage(sett, msg4.ToBytes())
	if err := handler(ctx, nil, nil, netMsg4); err != nil {
		t.Error(err) // works only logger
		return
	}

	netMsg5 := node.testNewNetworkMessage(sett, []byte{123})
	if err := handler(ctx, nil, nil, netMsg5); err == nil {
		t.Error("got success code with invalid message body")
		return
	}

	node.GetDBWrapper().Set(nil)
	netMsg41 := node.testNewNetworkMessage(sett, msg4.ToBytes())
	if err := handler(ctx, nil, nil, netMsg41); err == nil {
		t.Error("got success code with closed database")
		return
	}
}

func TestStoreHashWithBroadcastMessage(t *testing.T) {
	t.Parallel()

	_node, cancel := testNewNode(time.Minute, "", 6, 0, 0, false)
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
	netMsg := node.testNewNetworkMessage(sett, msg.ToBytes())
	logBuilder := anon_logger.NewLogBuilder("_")

	ctx := context.Background()
	if ok, err := node.storeHashWithBroadcast(ctx, logBuilder, netMsg); !ok || err != nil {
		t.Error(err)
		return
	}

	if ok, err := node.storeHashWithBroadcast(ctx, logBuilder, netMsg); ok || err != nil {
		switch {
		case ok:
			t.Error("success store one message again")
		case err != nil:
			t.Error("got error with try store twice same message")
		}
		return
	}

	db := node.GetDBWrapper().Get()
	node.GetDBWrapper().Set(nil)
	if ok, err := node.storeHashWithBroadcast(ctx, logBuilder, netMsg); ok || err == nil {
		t.Error("success use store function with null database")
		return
	}

	node.GetDBWrapper().Set(db)
	db.Close()
	if ok, err := node.storeHashWithBroadcast(ctx, logBuilder, netMsg); ok || err == nil {
		t.Error("success use store function with closed database")
		return
	}
}

func TestRecvSendMessage(t *testing.T) {
	t.Parallel()

	_node, cancel := testNewNode(time.Minute, "", 5, 0, 0, false)
	defer testFreeNodes([]INode{_node}, []context.CancelFunc{cancel}, 5)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	node := _node.(*sNode)
	if _, err := node.recvResponse(ctx, "not_exist"); err == nil {
		t.Error("success got action by undefined key")
		return
	}

	client := node.fQueue.GetClient()
	pubKey := client.GetPubKey()
	actionKey := newActionKey(pubKey, 111)

	node.setAction(actionKey)
	action, ok := node.getAction(actionKey)
	if !ok {
		t.Error("undefined created action key (1)")
		return
	}

	close(action)
	if _, err := node.recvResponse(ctx, actionKey); err == nil {
		t.Error("success got closed action")
		return
	}

	node.setAction(actionKey)
	if _, ok := node.getAction(actionKey); !ok {
		t.Error("undefined created action key (2)")
		return
	}

	cancel()
	if _, err := node.recvResponse(ctx, actionKey); err == nil {
		t.Error("success got action from canceled context")
		return
	}

	ctx2 := context.Background()

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
		if err := node.enqueueMessage(ctx2, msg); err != nil {
			t.Error("failed send message (push to queue)")
			return
		}
	}

	hasError := false
	for i := 0; i < 10; i++ {
		// message can be dequeued in the send's call time
		if err := node.enqueueMessage(ctx2, msg); err != nil {
			hasError = true
			break
		}
	}

	if !hasError {
		t.Error("success send message (push to queue) over queue capacity")
	}
}

func TestRetryEnqueue(t *testing.T) {
	t.Parallel()

	ctxBg, cancelBg := context.WithCancel(context.Background())
	defer cancelBg()

	_node, cancel := testNewNode(time.Minute, "", 11, 0, 3, false)
	defer testFreeNodes([]INode{_node}, []context.CancelFunc{cancel}, 11)

	node := _node.(*sNode)

	client := node.fQueue.GetClient()
	pubKey := client.GetPubKey()

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

	go func() {
		for i := 0; i < testutils.TCQueueCapacity*2; i++ {
			_ = node.enqueueMessage(ctxBg, msg)
		}
	}()

	time.Sleep(2 * time.Second)
	cancelBg()
	time.Sleep(time.Second)
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
		nodes[i], cancels[i] = testNewNode(timeWait, addrs[i], typeDB, i, 0, false)
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
			func(_ context.Context, _ INode, _ asymmetric.IPubKey, reqBytes []byte) ([]byte, error) {
				// send response
				return []byte(fmt.Sprintf("%s (response)", string(reqBytes))), nil
			},
		)
	}

	ctx := context.Background()
	go func() {
		if err := nodes[2].GetNetworkNode().Listen(ctx); err != nil && !errors.Is(err, net.ErrClosed) {
			t.Error(err)
		}
	}()
	go func() {
		if err := nodes[4].GetNetworkNode().Listen(ctx); err != nil && !errors.Is(err, net.ErrClosed) {
			t.Error(err)
		}
	}()

	// try connect to new node listeners
	// nodes to routes (nodes[0] -> nodes[2], nodes[1] -> nodes[4])
	err1 := testutils.TryN(50, 10*time.Millisecond, func() error {
		return nodes[0].GetNetworkNode().AddConnection(ctx, addresses[0])
	})
	if err1 != nil {
		t.Error(err1)
		return [5]INode{}, [5]context.CancelFunc{}
	}
	err2 := testutils.TryN(50, 10*time.Millisecond, func() error {
		return nodes[1].GetNetworkNode().AddConnection(ctx, addresses[1])
	})
	if err2 != nil {
		t.Error(err2)
		return [5]INode{}, [5]context.CancelFunc{}
	}

	// routes to routes (nodes[3] -> nodes[2], nodes[3] -> nodes[4])
	if err := nodes[3].GetNetworkNode().AddConnection(ctx, addresses[0]); err != nil {
		t.Error(err)
		return [5]INode{}, [5]context.CancelFunc{}
	}
	if err := nodes[3].GetNetworkNode().AddConnection(ctx, addresses[1]); err != nil {
		t.Error(err)
		return [5]INode{}, [5]context.CancelFunc{}
	}

	go func() {
		if err := nodes[0].Run(ctx); err == nil {
			t.Error("success twice running node")
			return
		}
	}()

	return nodes, cancels
}

/*
import (
	internal_anon_logger "github.com/number571/go-peer/internal/logger/anon"
	internal_std_logger "github.com/number571/go-peer/internal/logger/std"
)

type stLogging struct{}

func (p *stLogging) HasInfo() bool {
	return true
}
func (p *stLogging) HasWarn() bool {
	return true
}
func (p *stLogging) HasErro() bool {
	return true
}
*/

func testNewNode(timeWait time.Duration, addr string, typeDB, numDB, retryNum int, f2fDisabled bool) (INode, context.CancelFunc) {
	db, err := database.NewKVDatabase(
		database.NewSettings(&database.SSettings{
			FPath:     fmt.Sprintf(tcPathDBTemplate, typeDB, numDB),
			FWorkSize: testutils.TCWorkSize,
			FPassword: "CIPHER",
		}),
	)
	if err != nil {
		return nil, nil
	}
	parallel := uint64(1)
	networkMask := uint64(1)
	networkKey := "old_network_key"
	limitVoidSize := uint64(10_000)
	node := NewNode(
		NewSettings(&SSettings{
			FServiceName:  "TEST",
			FF2FDisabled:  f2fDisabled,
			FNetworkMask:  networkMask,
			FFetchTimeout: timeWait,
			FRetryEnqueue: uint64(retryNum),
		}),
		// internal_std_logger.NewStdLogger(&stLogging{}, internal_anon_logger.GetLogFunc()),
		logger.NewLogger(
			logger.NewSettings(&logger.SSettings{}),
			func(_ logger.ILogArg) string { return "" },
		),
		NewDBWrapper().Set(db),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:      addr,
				FMaxConnects:  testutils.TCMaxConnects,
				FReadTimeout:  timeWait,
				FWriteTimeout: timeWait,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FWorkSizeBits:          testutils.TCWorkSize,
					FLimitMessageSizeBytes: testutils.TCMessageSize,
					FLimitVoidSizeBytes:    limitVoidSize,
					FWaitReadTimeout:       time.Hour,
					FDialTimeout:           time.Minute,
					FReadTimeout:           time.Minute,
					FWriteTimeout:          time.Minute,
				}),
			}),
			conn.NewVSettings(&conn.SVSettings{
				FNetworkKey: networkKey,
			}),
			lru.NewLRUCache(
				lru.NewSettings(&lru.SSettings{
					FCapacity: testutils.TCCapacity,
				}),
			),
		),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{
				FNetworkMask:        networkMask,
				FWorkSizeBits:       testutils.TCWorkSize,
				FMainCapacity:       testutils.TCQueueCapacity,
				FVoidCapacity:       testutils.TCQueueCapacity,
				FParallel:           parallel,
				FLimitVoidSizeBytes: limitVoidSize,
				FDuration:           time.Second,
			}),
			queue.NewVSettings(&queue.SVSettings{
				FNetworkKey: networkKey,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FMessageSizeBytes: testutils.TCMessageSize,
					FKeySizeBits:      testutils.TcKeySize,
				}),
				asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
			),
		),
		asymmetric.NewListPubKeys(),
	)
	ctx, cancel := context.WithCancel(context.Background())
	go func() { _ = node.Run(ctx) }()
	return node, cancel
}

func testFreeNodes(nodes []INode, cancels []context.CancelFunc, typeDB int) {
	for i, node := range nodes {
		node.GetDBWrapper().Close()
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

func (p *sNode) testNewNetworkMessage(pSett net_message.ISettings, pMsgBytes []byte) net_message.IMessage {
	return net_message.NewMessage(
		pSett,
		payload.NewPayload(
			p.fSettings.GetNetworkMask(),
			pMsgBytes,
		),
		1,
		0,
	)
}
