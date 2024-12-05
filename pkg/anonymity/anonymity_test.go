// nolint: goerr113
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

	"github.com/number571/go-peer/pkg/anonymity/queue"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	testutils "github.com/number571/go-peer/test/utils"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/logger"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

const (
	tcPathDBTemplate = "database_test_%d_%d.db"
	tcIter           = 10
	tcWorkSize       = 10
	tcHead           = 123
	tcQueueCap       = 16
	tcMsgSize        = (8 << 10)
	tcMsgBody        = "hello, world!"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SAnonymityError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestNodeSettings(t *testing.T) {
	t.Parallel()

	node, networkNode, cancels := testNewNodeWithDB(time.Minute, "", &tsDatabase{})
	defer testFreeNodes([]INode{node}, []network.INode{networkNode}, []context.CancelFunc{cancels}, 9)

	sett := node.GetSettings()
	if sett.GetFetchTimeout() != time.Minute {
		t.Error("sett.GetFetchTimeout() != time.Minute")
		return
	}
	_ = node.GetLogger()

	_node := node.(*sNode)
	err := _node.storeHashIntoDatabase(
		anon_logger.NewLogBuilder("_"),
		hashing.NewHasher([]byte{}).ToBytes(),
	)
	if err == nil {
		t.Error("success store hash into database without correct set function")
		return
	}
}

func TestSettings(t *testing.T) {
	t.Parallel()

	for i := 0; i < 1; i++ {
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
	switch n { // nolint: gocritic
	case 0:
		_ = NewSettings(&SSettings{})
	}
}

func TestComplexFetchPayload(t *testing.T) {
	t.Parallel()

	addresses := [2]string{testutils.TgAddrs[2], testutils.TgAddrs[3]}
	nodes, networkNodes, cancels := testNewNodes(t, time.Minute, addresses, 0)
	if nodes[0] == nil {
		t.Error("nodes is null")
		return
	}
	defer testFreeNodes(nodes[:], networkNodes[:], cancels[:], 0)

	wg := sync.WaitGroup{}
	wg.Add(tcIter)

	ctx := context.Background()

	for i := 0; i < tcIter; i++ {
		go func(i int) {
			defer wg.Done()
			reqBody := fmt.Sprintf("%s (%d)", tcMsgBody, i)

			// nodes[1] -> nodes[0] -> nodes[2]
			resp, err := nodes[0].FetchPayload(
				ctx,
				nodes[1].GetQBProcessor().GetClient().GetPrivKey().GetPubKey(),
				payload.NewPayload32(tcHead, []byte(reqBody)),
			)
			if err != nil {
				t.Errorf("%s (%d)", err.Error(), i)
				return
			}

			if string(resp) != reqBody+" (response)" {
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
	addresses := [2]string{testutils.TgAddrs[10], testutils.TgAddrs[11]}
	nodes, networkNodes, cancels := testNewNodes(t, 3*time.Second, addresses, 1)
	if nodes[0] == nil {
		t.Error("nodes is null")
		return
	}
	defer testFreeNodes(nodes[:], networkNodes[:], cancels[:], 1)

	nodes[0].GetMapPubKeys().DelPubKey(nodes[1].GetQBProcessor().GetClient().GetPrivKey().GetPubKey())
	nodes[1].GetMapPubKeys().DelPubKey(nodes[0].GetQBProcessor().GetClient().GetPrivKey().GetPubKey())

	ctx := context.Background()

	// nodes[1] -> nodes[0] -> nodes[2]
	_, err := nodes[0].FetchPayload(
		ctx,
		nodes[1].GetQBProcessor().GetClient().GetPrivKey().GetPubKey(),
		payload.NewPayload32(tcHead, []byte(tcMsgBody)),
	)
	if err != nil {
		return
	}

	t.Error("get response without list of friends")
}

func TestFetchPayload(t *testing.T) {
	t.Parallel()

	addresses := [2]string{testutils.TgAddrs[12], testutils.TgAddrs[13]}
	nodes, networkNodes, cancels := testNewNodes(t, time.Minute, addresses, 4)
	if nodes[0] == nil {
		t.Error("nodes is null")
		return
	}
	defer testFreeNodes(nodes[:], networkNodes[:], cancels[:], 4)

	nodes[1].HandleFunc(
		tcHead,
		func(_ context.Context, _ INode, _ asymmetric.IPubKey, reqBytes []byte) ([]byte, error) {
			return []byte(fmt.Sprintf("echo: '%s'", string(reqBytes))), nil
		},
	)

	largeBodySize := nodes[0].GetQBProcessor().GetClient().GetPayloadLimit() - encoding.CSizeUint64 + 1
	ctx := context.Background()
	_, err := nodes[0].FetchPayload(
		ctx,
		nodes[1].GetQBProcessor().GetClient().GetPrivKey().GetPubKey(),
		payload.NewPayload32(tcHead, random.NewRandom().GetBytes(largeBodySize)),
	)
	if err == nil {
		t.Error("success fetch payload with large body")
		return
	}

	result, err1 := nodes[0].FetchPayload(
		ctx,
		nodes[1].GetQBProcessor().GetClient().GetPrivKey().GetPubKey(),
		payload.NewPayload32(tcHead, []byte(tcMsgBody)),
	)
	if err1 != nil {
		t.Error(err1)
		return
	}

	if string(result) != fmt.Sprintf("echo: '%s'", tcMsgBody) {
		t.Error("got invalid message body")
		return
	}
}

func TestBroadcastPayload(t *testing.T) {
	t.Parallel()

	addresses := [2]string{testutils.TgAddrs[14], testutils.TgAddrs[15]}
	nodes, networkNodes, cancels := testNewNodes(t, time.Minute, addresses, 3)
	if nodes[0] == nil {
		t.Error("nodes is null")
		return
	}
	defer testFreeNodes(nodes[:], networkNodes[:], cancels[:], 3)

	chResult := make(chan string)
	nodes[1].HandleFunc(
		tcHead,
		func(_ context.Context, _ INode, _ asymmetric.IPubKey, reqBytes []byte) ([]byte, error) {
			res := fmt.Sprintf("echo: '%s'", string(reqBytes))
			go func() { chResult <- res }()
			return nil, nil
		},
	)

	largeBodySize := nodes[0].GetQBProcessor().GetClient().GetPayloadLimit() - encoding.CSizeUint64 + 1
	err := nodes[0].SendPayload(
		context.Background(),
		nodes[1].GetQBProcessor().GetClient().GetPrivKey().GetPubKey(),
		payload.NewPayload64(uint64(tcHead), random.NewRandom().GetBytes(largeBodySize)),
	)
	if err == nil {
		t.Error("success broadcast payload with large body")
		return
	}

	err1 := nodes[0].SendPayload(
		context.Background(),
		nodes[1].GetQBProcessor().GetClient().GetPrivKey().GetPubKey(),
		payload.NewPayload64(uint64(tcHead), []byte(tcMsgBody)),
	)
	if err1 != nil {
		t.Error(err1)
		return
	}

	select {
	case x := <-chResult:
		if x != fmt.Sprintf("echo: '%s'", tcMsgBody) {
			t.Error("got invalid message body")
			return
		}
		// success
	case <-time.After(time.Minute):
		t.Error("error: time after 1 minute")
		return
	}
}

func TestEnqueuePayload(t *testing.T) {
	t.Parallel()

	addresses := [2]string{testutils.TgAddrs[16], testutils.TgAddrs[17]}
	nodes, networkNodes, cancels := testNewNodes(t, time.Minute, addresses, 8)
	if nodes[0] == nil {
		t.Error("nodes is null")
		return
	}
	defer testFreeNodes(nodes[:], networkNodes[:], cancels[:], 8)

	node := nodes[0].(*sNode)
	pubKey := nodes[1].GetQBProcessor().GetClient().GetPrivKey().GetPubKey()

	logBuilder := anon_logger.NewLogBuilder("test")
	pld := payload.NewPayload64(uint64(tcHead), []byte(tcMsgBody))

	overheadBody := random.NewRandom().GetBytes(tcMsgSize + 1)
	overPld := payload.NewPayload64(uint64(tcHead), overheadBody)
	if err := node.enqueuePayload(logBuilder, pubKey, overPld); err == nil {
		t.Error("success with overhead message")
		return
	}

	pldBytes := payload.NewPayload64(
		joinHead(sAction(1).setType(true), tcHead).uint64(),
		[]byte(tcMsgBody),
	).ToBytes()

	for i := 0; i < tcQueueCap; i++ {
		if err := node.fQBProcessor.EnqueueMessage(pubKey, pldBytes); err != nil {
			t.Error("failed send message (push to queue)")
			return
		}
	}

	// after full queue
	for i := 0; i < 2*tcQueueCap; i++ {
		if err := node.enqueuePayload(logBuilder, pubKey, pld); err != nil {
			return
		}
	}

	t.Error("success enqueue payload over queue capacity")
}

func TestHandleWrapper(t *testing.T) {
	t.Parallel()

	_node, networkNode, cancel := testNewNode(time.Minute, "", 7, 0)
	defer testFreeNodes([]INode{_node}, []network.INode{networkNode}, []context.CancelFunc{cancel}, 7)

	node := _node.(*sNode)
	handler := node.messageHandler
	client := node.fQBProcessor.GetClient()

	privKey := client.GetPrivKey()
	pubKey := privKey.GetPubKey()
	node.GetMapPubKeys().SetPubKey(privKey.GetPubKey())

	ctx := context.Background()
	sett := net_message.NewConstructSettings(&net_message.SConstructSettings{
		FSettings: net_message.NewSettings(&net_message.SSettings{}),
	})

	msg, err := client.EncryptMessage(
		pubKey,
		payload.NewPayload64(
			joinHead(sAction(1).setType(true), tcHead).uint64(),
			[]byte(tcMsgBody),
		).ToBytes(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	netMsg := node.testNewNetworkMessage(sett, msg)
	if err := handler(ctx, netMsg); err != nil {
		t.Error(err)
		return
	}

	if err := handler(ctx, netMsg); err != nil {
		t.Error("repeated message:", err.Error())
		return
	}

	msgWithoutPld, err := client.EncryptMessage(pubKey, []byte{123})
	if err != nil {
		t.Error(err)
		return
	}

	netMsgWithoutPld := node.testNewNetworkMessage(sett, msgWithoutPld)
	if err := handler(ctx, netMsgWithoutPld); err != nil {
		t.Error(err) // works only logger
		return
	}

	node.HandleFunc(
		111,
		func(_ context.Context, _ INode, _ asymmetric.IPubKey, _ []byte) ([]byte, error) {
			return nil, errors.New("some error")
		},
	)

	msg2, err := client.EncryptMessage(
		pubKey,
		payload.NewPayload64(
			joinHead(sAction(1).setType(true), 111).uint64(),
			[]byte(tcMsgBody),
		).ToBytes(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	netMsg2 := node.testNewNetworkMessage(sett, msg2)
	if err := handler(ctx, netMsg2); err != nil {
		t.Error(err) // works only logger
		return
	}

	msg3, err := client.EncryptMessage(
		pubKey,
		payload.NewPayload64(
			uint64(111),
			[]byte("?"+tcMsgBody),
		).ToBytes(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	netMsg3 := node.testNewNetworkMessage(sett, msg3)
	if err := handler(ctx, netMsg3); err != nil {
		t.Error(err) // works only logger
		return
	}

	msg4, err := client.EncryptMessage(
		pubKey,
		payload.NewPayload64(
			joinHead(sAction(1).setType(false), 111).uint64(),
			[]byte(tcMsgBody),
		).ToBytes(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	netMsg4 := node.testNewNetworkMessage(sett, msg4)
	if err := handler(ctx, netMsg4); err != nil {
		t.Error(err) // works only logger
		return
	}

	netMsg5 := node.testNewNetworkMessage(sett, []byte{123})
	if err := handler(ctx, netMsg5); err == nil {
		t.Error("got success code with invalid message body")
		return
	}

	node.fKVDatavase.Close()
	netMsg41 := node.testNewNetworkMessage(sett, msg4)
	if err := handler(ctx, netMsg41); err == nil {
		t.Error("got success code with closed database")
		return
	}
}

func TestStoreHashWithBroadcastMessage(t *testing.T) {
	t.Parallel()

	_node, networkNode, cancel := testNewNode(time.Minute, "", 6, 0)
	defer testFreeNodes([]INode{_node}, []network.INode{networkNode}, []context.CancelFunc{cancel}, 6)

	node := _node.(*sNode)
	client := node.fQBProcessor.GetClient()

	msg, err := client.EncryptMessage(
		client.GetPrivKey().GetPubKey(),
		payload.NewPayload64(
			joinHead(sAction(1).setType(true), 111).uint64(),
			[]byte(tcMsgBody),
		).ToBytes(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	sett := net_message.NewConstructSettings(&net_message.SConstructSettings{
		FSettings: net_message.NewSettings(&net_message.SSettings{}),
	})

	netMsg := node.testNewNetworkMessage(sett, msg)
	logBuilder := anon_logger.NewLogBuilder("_")

	ctx := context.Background()
	if ok, err := node.storeHashWithProduce(ctx, logBuilder, netMsg); !ok || err != nil {
		t.Error(err)
		return
	}

	if ok, err := node.storeHashWithProduce(ctx, logBuilder, netMsg); ok || err != nil {
		switch {
		case ok:
			t.Error("success store one message again")
		case err != nil:
			t.Error("got error with try store twice same message")
		}
		return
	}
}

func TestRecvSendMessage(t *testing.T) {
	t.Parallel()

	_node, networkNode, cancel := testNewNode(time.Minute, "", 5, 0)
	defer testFreeNodes([]INode{_node}, []network.INode{networkNode}, []context.CancelFunc{cancel}, 5)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	node := _node.(*sNode)
	if _, err := node.recvResponse(ctx, "not_exist"); err == nil {
		t.Error("success got action by undefined key")
		return
	}

	client := node.fQBProcessor.GetClient()
	pubKey := client.GetPrivKey().GetPubKey()
	actionKey := newActionKey(pubKey, sAction(111).setType(true))

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

	msgBody := "hello, world!"
	pldBytes := payload.NewPayload64(
		joinHead(sAction(1).setType(true), tcHead).uint64(),
		[]byte(msgBody),
	).ToBytes()

	for i := 0; i < tcQueueCap; i++ {
		if err := node.fQBProcessor.EnqueueMessage(pubKey, pldBytes); err != nil {
			t.Error("failed send message (push to queue)")
			return
		}
	}

	hasError := false
	for i := 0; i < 10; i++ {
		// message can be dequeued in the send's call time
		if err := node.fQBProcessor.EnqueueMessage(pubKey, pldBytes); err != nil {
			hasError = true
			break
		}
	}

	if !hasError {
		t.Error("success send message (push to queue) over queue capacity")
	}
}

// nodes[0], nodes[1] = clients
// nodes[2], nodes[3], nodes[4] = routes
// nodes[2], nodes[4] are have open ports
// Scheme: (nodes[0]) -> nodes[2] -> nodes[3] -> nodes[4] -> (nodes[1])
func testNewNodes(t *testing.T, timeWait time.Duration, addresses [2]string, typeDB int) ([5]INode, [5]network.INode, [5]context.CancelFunc) {
	nodes := [5]INode{}
	networkNodes := [5]network.INode{}
	cancels := [5]context.CancelFunc{}
	addrs := [5]string{"", "", addresses[0], "", addresses[1]}

	for i := 0; i < 5; i++ {
		nodes[i], networkNodes[i], cancels[i] = testNewNode(timeWait, addrs[i], typeDB, i)
		if nodes[i] == nil {
			t.Errorf("node (%d) is not running %d", i, typeDB)
			return [5]INode{}, [5]network.INode{}, [5]context.CancelFunc{}
		}
	}

	pubKey1 := nodes[1].GetQBProcessor().GetClient().GetPrivKey().GetPubKey()
	pubKey0 := nodes[0].GetQBProcessor().GetClient().GetPrivKey().GetPubKey()

	nodes[0].GetMapPubKeys().SetPubKey(pubKey1)
	nodes[1].GetMapPubKeys().SetPubKey(pubKey0)

	for _, node := range nodes {
		node.HandleFunc(
			tcHead,
			func(_ context.Context, _ INode, _ asymmetric.IPubKey, reqBytes []byte) ([]byte, error) {
				// send response
				return []byte(string(reqBytes) + " (response)"), nil
			},
		)
	}

	ctx := context.Background()
	go func() {
		if err := networkNodes[2].Listen(ctx); err != nil && !errors.Is(err, net.ErrClosed) {
			t.Error(err)
		}
	}()
	go func() {
		if err := networkNodes[4].Listen(ctx); err != nil && !errors.Is(err, net.ErrClosed) {
			t.Error(err)
		}
	}()

	// try connect to new node listeners
	// nodes to routes (nodes[0] -> nodes[2], nodes[1] -> nodes[4])
	err1 := testutils.TryN(50, 10*time.Millisecond, func() error {
		return networkNodes[0].AddConnection(ctx, addresses[0])
	})
	if err1 != nil {
		t.Error(err1)
		return [5]INode{}, [5]network.INode{}, [5]context.CancelFunc{}
	}
	err2 := testutils.TryN(50, 10*time.Millisecond, func() error {
		return networkNodes[1].AddConnection(ctx, addresses[1])
	})
	if err2 != nil {
		t.Error(err2)
		return [5]INode{}, [5]network.INode{}, [5]context.CancelFunc{}
	}

	// routes to routes (nodes[3] -> nodes[2], nodes[3] -> nodes[4])
	if err := networkNodes[3].AddConnection(ctx, addresses[0]); err != nil {
		t.Error(err)
		return [5]INode{}, [5]network.INode{}, [5]context.CancelFunc{}
	}
	if err := networkNodes[3].AddConnection(ctx, addresses[1]); err != nil {
		t.Error(err)
		return [5]INode{}, [5]network.INode{}, [5]context.CancelFunc{}
	}

	go func() {
		if err := nodes[0].Run(ctx); err == nil {
			t.Error("success twice running node")
			return
		}
	}()

	return nodes, networkNodes, cancels
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

func testNewNodeWithDB(timeWait time.Duration, addr string, db database.IKVDatabase) (INode, network.INode, context.CancelFunc) {
	msgChan := make(chan net_message.IMessage)
	parallel := uint64(1)
	networkMask := uint32(1)
	limitVoidSize := uint64(10_000)
	networkNode := network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      addr,
			FMaxConnects:  16,
			FReadTimeout:  timeWait,
			FWriteTimeout: timeWait,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSettings: net_message.NewSettings(&net_message.SSettings{
					FWorkSizeBits: tcWorkSize,
				}),
				FLimitMessageSizeBytes: tcMsgSize + limitVoidSize,
				FWaitReadTimeout:       time.Hour,
				FDialTimeout:           time.Minute,
				FReadTimeout:           time.Minute,
				FWriteTimeout:          time.Minute,
			}),
		}),
		cache.NewLRUCache(1024),
	).HandleFunc(networkMask, func(_ context.Context, _ network.INode, _ conn.IConn, msg net_message.IMessage) error {
		msgChan <- msg
		return nil
	})
	node := NewNode(
		NewSettings(&SSettings{
			FServiceName:  "TEST",
			FFetchTimeout: timeWait,
		}),
		// internal_std_logger.NewStdLogger(&stLogging{}, internal_anon_logger.GetLogFunc()),
		logger.NewLogger(
			logger.NewSettings(&logger.SSettings{}),
			func(_ logger.ILogArg) string { return "" },
		),
		NewAdapterByFuncs(
			func(ctx context.Context, msg net_message.IMessage) error {
				return networkNode.BroadcastMessage(ctx, msg)
			},
			func(ctx context.Context) (net_message.IMessage, error) {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case msg := <-msgChan:
					return msg, nil
				}
			},
		),
		db,
		queue.NewQBProblemProcessor(
			queue.NewSettings(&queue.SSettings{
				FMessageConstructSettings: net_message.NewConstructSettings(&net_message.SConstructSettings{
					FSettings: net_message.NewSettings(&net_message.SSettings{
						FWorkSizeBits: tcWorkSize,
					}),
					FParallel:             parallel,
					FRandMessageSizeBytes: limitVoidSize,
				}),
				FNetworkMask:  networkMask,
				FQueuePoolCap: [2]uint64{tcQueueCap, tcQueueCap},
				FQueuePeriod:  time.Second,
				FConsumersCap: 1,
			}),
			client.NewClient(
				asymmetric.NewPrivKey(),
				tcMsgSize,
			),
		),
	)
	ctx, cancel := context.WithCancel(context.Background())
	go func() { _ = node.Run(ctx) }()
	return node, networkNode, cancel
}

func testNewNode(timeWait time.Duration, addr string, typeDB, numDB int) (INode, network.INode, context.CancelFunc) {
	db, err := database.NewKVDatabase(fmt.Sprintf(tcPathDBTemplate, typeDB, numDB))
	if err != nil {
		panic(err)
	}
	return testNewNodeWithDB(timeWait, addr, db)
}

func testFreeNodes(nodes []INode, networks []network.INode, cancels []context.CancelFunc, typeDB int) {
	for i, node := range nodes {
		node.GetKVDatabase().Close()
		networks[i].Close()
		cancels[i]()
	}
	testDeleteDB(typeDB)
}

func testDeleteDB(typeDB int) {
	for i := 0; i < 5; i++ {
		os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, typeDB, i))
	}
}

func (p *sNode) testNewNetworkMessage(pSett net_message.IConstructSettings, pMsgBytes []byte) net_message.IMessage {
	return net_message.NewMessage(
		pSett,
		payload.NewPayload32(
			p.fQBProcessor.GetSettings().GetNetworkMask(),
			pMsgBytes,
		),
	)
}

type tsDatabase struct{}

func (p *tsDatabase) Get([]byte) ([]byte, error) { return nil, database.ErrNotFound }
func (p *tsDatabase) Set([]byte, []byte) error   { return errors.New("some error") }
func (p *tsDatabase) Del([]byte) error           { return nil }
func (p *tsDatabase) Close() error               { return nil }
