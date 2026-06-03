// nolint: err113
package qb

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/anonymity/qb/adapters"
	"github.com/number571/go-peer/pkg/anonymity/qb/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2/hybrid"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	testutils "github.com/number571/go-peer/test/utils"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/qb/logger"
	"github.com/number571/go-peer/pkg/network/conn"
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
		t.Fatal("incorrect err.Error()")
	}
}

func TestNodeSettings(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	node, _, _ := testRunNodeWithDB(ctx, time.Minute, "", &tsDatabase{}, func(ctx context.Context, i INode, ik layer2.IParticipantKey, b []byte) ([]byte, error) {
		return nil, nil
	})
	defer testFreeNodes([]INode{node}, 9)

	sett := node.GetSettings()
	if sett.GetFetchTimeout() != time.Minute {
		t.Fatal("sett.GetFetchTimeout() != time.Minute")
	}
	_ = node.GetLogger()

	_node := node.(*sNode)
	err := _node.storeHashIntoDatabase(
		anon_logger.NewLogBuilder("_"),
		layer1.NewMessage(
			layer1.NewConstructSettings(&layer1.SConstructSettings{
				FSettings: layer1.NewSettings(&layer1.SSettings{}),
			}),
			[]byte{},
		),
	)
	if err == nil {
		t.Fatal("success store hash into database without correct set function")
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
			t.Fatal("nothing panics")
		}
	}()
	switch n { // nolint: gocritic
	case 0:
		_ = NewSettings(&SSettings{})
	}
}

func TestComplexFetchPayload(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addresses := [2]string{testutils.TgAddrs[2], testutils.TgAddrs[3]}
	nodes, privKeys := testRunNodes(ctx, t, time.Minute, addresses, 0, func(_ context.Context, _ INode, _ layer2.IParticipantKey, reqBytes []byte) ([]byte, error) {
		return []byte(string(reqBytes) + " (response)"), nil
	})
	if nodes[0] == nil {
		t.Fatal("nodes is null")
	}
	defer testFreeNodes(nodes[:], 0)

	wg := sync.WaitGroup{}
	wg.Add(tcIter)

	for i := 0; i < tcIter; i++ {
		go func(i int) {
			defer wg.Done()
			reqBody := fmt.Sprintf("%s (%d)", tcMsgBody, i)

			// nodes[1] -> nodes[0] -> nodes[2]
			resp, err := nodes[0].FetchPayload(
				ctx,
				privKeys[1].GetPubKey(),
				[]byte(reqBody),
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 3 seconds for wait
	addresses := [2]string{testutils.TgAddrs[10], testutils.TgAddrs[11]}
	nodes, privKeys := testRunNodes(ctx, t, 3*time.Second, addresses, 1, func(_ context.Context, _ INode, _ layer2.IParticipantKey, reqBytes []byte) ([]byte, error) {
		return []byte(string(reqBytes) + " (response)"), nil
	})
	if nodes[0] == nil {
		t.Fatal("nodes is null")
	}
	defer testFreeNodes(nodes[:], 1)

	nodes[0].GetKeysContainer().Del(privKeys[1].GetPubKey().GetHasher().ToString())
	nodes[1].GetKeysContainer().Del(privKeys[0].GetPubKey().GetHasher().ToString())

	// nodes[1] -> nodes[0] -> nodes[2]
	_, err := nodes[0].FetchPayload(
		ctx,
		privKeys[1].GetPubKey(),
		[]byte(tcMsgBody),
	)
	if err != nil {
		return
	}

	t.Fatal("get response without list of friends")
}

func TestFetchPayload(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addresses := [2]string{testutils.TgAddrs[12], testutils.TgAddrs[13]}
	nodes, privKeys := testRunNodes(ctx, t, time.Minute, addresses, 4, func(_ context.Context, _ INode, _ layer2.IParticipantKey, reqBytes []byte) ([]byte, error) {
		return []byte(fmt.Sprintf("echo: '%s'", string(reqBytes))), nil
	})
	if nodes[0] == nil {
		t.Fatal("nodes is null")
	}
	defer testFreeNodes(nodes[:], 4)

	largeBodySize := nodes[0].GetQBProcessor().GetScheme().GetPayloadLimit() - encoding.CSizeUint64 + 1
	_, err := nodes[0].FetchPayload(
		ctx,
		privKeys[1].GetPubKey(),
		random.NewRandom().GetBytes(largeBodySize),
	)
	if err == nil {
		t.Fatal("success fetch payload with large body")
	}

	result, err1 := nodes[0].FetchPayload(
		ctx,
		privKeys[1].GetPubKey(),
		[]byte(tcMsgBody),
	)
	if err1 != nil {
		t.Fatal(err1)
	}

	if string(result) != fmt.Sprintf("echo: '%s'", tcMsgBody) {
		t.Fatal("got invalid message body")
	}
}

func TestBroadcastPayload(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chResult := make(chan string)
	addresses := [2]string{testutils.TgAddrs[14], testutils.TgAddrs[15]}
	nodes, privKeys := testRunNodes(ctx, t, time.Minute, addresses, 3, func(_ context.Context, _ INode, _ layer2.IParticipantKey, reqBytes []byte) ([]byte, error) {
		res := fmt.Sprintf("echo: '%s'", string(reqBytes))
		go func() { chResult <- res }()
		return nil, nil
	})
	if nodes[0] == nil {
		t.Fatal("nodes is null")
	}
	defer testFreeNodes(nodes[:], 3)

	largeBodySize := nodes[0].GetQBProcessor().GetScheme().GetPayloadLimit() - encoding.CSizeUint64 + 1
	err := nodes[0].SendPayload(
		context.Background(),
		privKeys[1].GetPubKey(),
		random.NewRandom().GetBytes(largeBodySize),
	)
	if err == nil {
		t.Fatal("success broadcast payload with large body")
	}

	err1 := nodes[0].SendPayload(
		context.Background(),
		privKeys[1].GetPubKey(),
		[]byte(tcMsgBody),
	)
	if err1 != nil {
		t.Fatal(err1)
	}

	select {
	case x := <-chResult:
		if x != fmt.Sprintf("echo: '%s'", tcMsgBody) {
			t.Fatal("got invalid message body")
		}
		// success
	case <-time.After(time.Minute):
		t.Fatal("error: time after 1 minute")
	}
}

func TestEnqueuePayload(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addresses := [2]string{testutils.TgAddrs[16], testutils.TgAddrs[17]}
	nodes, privKeys := testRunNodes(ctx, t, time.Minute, addresses, 8, func(_ context.Context, _ INode, _ layer2.IParticipantKey, reqBytes []byte) ([]byte, error) {
		return []byte(string(reqBytes) + " (response)"), nil
	})
	if nodes[0] == nil {
		t.Fatal("nodes is null")
	}
	defer testFreeNodes(nodes[:], 8)

	node := nodes[0].(*sNode)
	pubKey := privKeys[1].GetPubKey()

	logBuilder := anon_logger.NewLogBuilder("test")
	pld := payload.NewPayload64(uint64(tcHead), []byte(tcMsgBody))

	overheadBody := random.NewRandom().GetBytes(tcMsgSize + 1)
	overPld := payload.NewPayload64(uint64(tcHead), overheadBody)
	if err := node.enqueuePayload(logBuilder, pubKey, overPld); err == nil {
		t.Fatal("success with overhead message")
	}

	pldBytes := payload.NewPayload64(
		setRequestBit(tcHead),
		[]byte(tcMsgBody),
	).ToBytes()

	for i := 0; i < tcQueueCap; i++ {
		if err := node.fQBProcessor.EnqueueMessage(pubKey, pldBytes); err != nil {
			t.Fatal("failed send message (push to queue)")
		}
	}

	// after full queue
	for i := 0; i < 2*tcQueueCap; i++ {
		if err := node.enqueuePayload(logBuilder, pubKey, pld); err != nil {
			return
		}
	}

	t.Fatal("success enqueue payload over queue capacity")
}

func TestHandleWrapper(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_node, _, privKey := testRunNode(ctx, time.Minute, "", 7, 0, func(_ context.Context, _ INode, _ layer2.IParticipantKey, _ []byte) ([]byte, error) {
		return nil, errors.New("some error") //nolint:err113
	})
	defer testFreeNodes([]INode{_node}, 7)

	node := _node.(*sNode)
	handler := node.consumeMessage
	scheme := node.fQBProcessor.GetScheme()

	pubKey := privKey.GetPubKey()
	node.GetKeysContainer().Add(privKey.GetPubKey())

	sett := layer1.NewConstructSettings(&layer1.SConstructSettings{
		FSettings: layer1.NewSettings(&layer1.SSettings{
			FWorkSizeBits: tcWorkSize,
		}),
	})

	msg, err := scheme.EncryptMessage(
		pubKey,
		payload.NewPayload64(
			setRequestBit(tcHead),
			[]byte(tcMsgBody),
		).ToBytes(),
	)
	if err != nil {
		t.Fatal(err)
	}

	netMsg := node.testNewNetworkMessage(sett, msg)
	if err := handler(ctx, netMsg); err != nil {
		t.Fatal(err)
	}

	netMsgY := node.testNewNetworkMessageWithAnotherMessageType(sett, msg)
	if err := handler(ctx, netMsgY); err == nil {
		t.Fatal("success handle message with another network type")
	}

	if err := handler(ctx, netMsg); err != nil {
		t.Fatal("repeated message:", err.Error())
	}

	msgWithoutPld, err := scheme.EncryptMessage(pubKey, []byte{123})
	if err != nil {
		t.Fatal(err)
	}

	netMsgWithoutPld := node.testNewNetworkMessage(sett, msgWithoutPld)
	if err := handler(ctx, netMsgWithoutPld); err != nil {
		t.Fatal(err) // works only logger
	}

	msg2, err := scheme.EncryptMessage(
		pubKey,
		payload.NewPayload64(
			setRequestBit(111),
			[]byte(tcMsgBody),
		).ToBytes(),
	)
	if err != nil {
		t.Fatal(err)
	}

	netMsg2 := node.testNewNetworkMessage(sett, msg2)
	if err := handler(ctx, netMsg2); err != nil {
		t.Fatal(err) // works only logger
	}

	msg3, err := scheme.EncryptMessage(
		pubKey,
		payload.NewPayload64(
			uint64(111),
			[]byte("?"+tcMsgBody),
		).ToBytes(),
	)
	if err != nil {
		t.Fatal(err)
	}

	netMsg3 := node.testNewNetworkMessage(sett, msg3)
	if err := handler(ctx, netMsg3); err != nil {
		t.Fatal(err) // works only logger
	}

	msg4, err := scheme.EncryptMessage(
		pubKey,
		payload.NewPayload64(
			setResponseBit(111),
			[]byte(tcMsgBody),
		).ToBytes(),
	)
	if err != nil {
		t.Fatal(err)
	}

	netMsg4 := node.testNewNetworkMessage(sett, msg4)
	if err := handler(ctx, netMsg4); err != nil {
		t.Fatal(err) // works only logger
	}

	netMsg5 := node.testNewNetworkMessage(sett, []byte{123})
	if err := handler(ctx, netMsg5); err == nil {
		t.Fatal("got success code with invalid message body")
	}

	_ = node.fKVDatavase.Close()
	netMsg41 := node.testNewNetworkMessage(sett, msg4)
	if err := handler(ctx, netMsg41); err == nil {
		t.Fatal("got success code with closed database")
	}
}

func TestStoreHashWithBroadcastMessage(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_node, _, privKey := testRunNode(ctx, time.Minute, "", 6, 0, func(ctx context.Context, i INode, ik layer2.IParticipantKey, b []byte) ([]byte, error) {
		return nil, nil
	})
	defer testFreeNodes([]INode{_node}, 6)

	node := _node.(*sNode)
	scheme := node.fQBProcessor.GetScheme()

	msg, err := scheme.EncryptMessage(
		privKey.GetPubKey(),
		payload.NewPayload64(
			setRequestBit(111),
			[]byte(tcMsgBody),
		).ToBytes(),
	)
	if err != nil {
		t.Fatal(err)
	}

	sett := layer1.NewConstructSettings(&layer1.SConstructSettings{
		FSettings: layer1.NewSettings(&layer1.SSettings{}),
	})

	netMsg := node.testNewNetworkMessage(sett, msg)
	if ok, err := node.produceMessage(ctx, netMsg); !ok || err != nil {
		t.Fatal(err)
	}
	if ok, err := node.produceMessage(ctx, netMsg); ok || err != nil {
		switch {
		case ok:
			t.Fatal("success store one message again")
		case err != nil:
			t.Fatal("got error with try store twice same message")
		}
		return
	}
}

func TestRecvSendMessage(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_node, _, privKey := testRunNode(ctx, time.Minute, "", 5, 0, func(ctx context.Context, i INode, ik layer2.IParticipantKey, b []byte) ([]byte, error) {
		return nil, nil
	})
	defer testFreeNodes([]INode{_node}, 5)

	node := _node.(*sNode)
	if _, err := node.recvResponse(ctx, "not_exist"); err == nil {
		t.Fatal("success got action by undefined key")
	}

	pubKey := privKey.GetPubKey()
	actionKey := newActionKey(pubKey, setRequestBit(111))

	node.setAction(actionKey)
	action, ok := node.getAction(actionKey)
	if !ok {
		t.Fatal("undefined created action key (1)")
	}

	close(action)
	if _, err := node.recvResponse(ctx, actionKey); err == nil {
		t.Fatal("success got closed action")
	}

	node.setAction(actionKey)
	if _, ok := node.getAction(actionKey); !ok {
		t.Fatal("undefined created action key (2)")
	}

	cancel()
	if _, err := node.recvResponse(ctx, actionKey); err == nil {
		t.Fatal("success got action from canceled context")
	}

	msgBody := "hello, world!"
	pldBytes := payload.NewPayload64(
		setRequestBit(tcHead),
		[]byte(msgBody),
	).ToBytes()

	for i := 0; i < tcQueueCap; i++ {
		if err := node.fQBProcessor.EnqueueMessage(pubKey, pldBytes); err != nil {
			t.Fatal("failed send message (push to queue)")
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
		t.Fatal("success send message (push to queue) over queue capacity")
	}
}

// nodes[0], nodes[1] = clients
// nodes[2], nodes[3], nodes[4] = routes
// nodes[2], nodes[4] are have open ports
// Scheme: (nodes[0]) -> nodes[2] -> nodes[3] -> nodes[4] -> (nodes[1])
func testRunNodes(ctx context.Context, t *testing.T, timeWait time.Duration, addresses [2]string, typeDB int, handlerF IHandlerF) ([5]INode, [5]asymmetric.IPrivKey) {
	nodes := [5]INode{}
	networkNodes := [5]network.INode{}
	privKeys := [5]asymmetric.IPrivKey{}
	addrs := [5]string{"", "", addresses[0], "", addresses[1]}

	for i := 0; i < 5; i++ {
		nodes[i], networkNodes[i], privKeys[i] = testRunNode(ctx, timeWait, addrs[i], typeDB, i, handlerF)
		if nodes[i] == nil {
			t.Errorf("node (%d) is not running %d", i, typeDB)
			return [5]INode{}, [5]asymmetric.IPrivKey{}
		}
	}

	pubKey1 := privKeys[1].GetPubKey()
	pubKey0 := privKeys[0].GetPubKey()

	nodes[0].GetKeysContainer().Add(pubKey1)
	nodes[1].GetKeysContainer().Add(pubKey0)

	go func() {
		if err := networkNodes[2].Run(ctx); err != nil && !errors.Is(err, net.ErrClosed) && !errors.Is(err, context.Canceled) {
			t.Error(err)
		}
	}()
	go func() {
		if err := networkNodes[4].Run(ctx); err != nil && !errors.Is(err, net.ErrClosed) && !errors.Is(err, context.Canceled) {
			t.Error(err)
		}
	}()

	// try connect to new node listeners
	// nodes to routes (nodes[0] -> nodes[2], nodes[1] -> nodes[4])
	err1 := testutils.TryN(50, 10*time.Millisecond, func() error {
		return networkNodes[0].AddConnection(ctx, addresses[0])
	})
	if err1 != nil {
		t.Fatal(err1)
		return [5]INode{}, [5]asymmetric.IPrivKey{}
	}
	err2 := testutils.TryN(50, 10*time.Millisecond, func() error {
		return networkNodes[1].AddConnection(ctx, addresses[1])
	})
	if err2 != nil {
		t.Fatal(err2)
		return [5]INode{}, [5]asymmetric.IPrivKey{}
	}

	// routes to routes (nodes[3] -> nodes[2], nodes[3] -> nodes[4])
	if err := networkNodes[3].AddConnection(ctx, addresses[0]); err != nil {
		t.Fatal(err)
		return [5]INode{}, [5]asymmetric.IPrivKey{}
	}
	if err := networkNodes[3].AddConnection(ctx, addresses[1]); err != nil {
		t.Fatal(err)
		return [5]INode{}, [5]asymmetric.IPrivKey{}
	}

	go func() {
		if err := nodes[0].Run(ctx); err == nil {
			t.Error("success twice running node")
		}
	}()

	return nodes, privKeys
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

func testRunNodeWithDB(ctx context.Context, timeWait time.Duration, addr string, db database.IKVDatabase, handlerF IHandlerF) (INode, network.INode, asymmetric.IPrivKey) {
	privKey := asymmetric.NewPrivKey()
	msgChan := make(chan layer1.IMessage)
	parallel := uint64(1)
	limitVoidSize := uint64(10_000)
	networkNode := network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      addr,
			FMaxConnects:  16,
			FReadTimeout:  timeWait,
			FWriteTimeout: timeWait,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSettings: layer1.NewSettings(&layer1.SSettings{
					FWorkSizeBits: tcWorkSize,
				}),
				FLimitMessageSizeBytes: tcMsgSize + limitVoidSize,
				FWaitReadTimeout:       time.Hour,
				FDialTimeout:           time.Minute,
				FReadTimeout:           time.Minute,
				FWriteTimeout:          time.Minute,
			}),
		}),
		func(_ context.Context, _ network.INode, _ conn.IConn, msg layer1.IMessage) error {
			msgChan <- msg
			return nil
		},
		cache.NewLRUCache(1024),
	)
	node := NewNode(
		NewSettings(&SSettings{
			FServiceName:  "TEST",
			FFetchTimeout: timeWait,
		}),
		handlerF,
		// internal_std_logger.NewStdLogger(&stLogging{}, internal_anon_logger.GetLogFunc()),
		logger.NewLogger(
			logger.NewSettings(&logger.SSettings{}),
			func(_ logger.ILogArg) string { return "" },
		),
		adapters.NewAdapterByFuncs(
			func(ctx context.Context, msg layer1.IMessage) error {
				return networkNode.BroadcastMessage(ctx, msg)
			},
			func(ctx context.Context) (layer1.IMessage, error) {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case msg := <-msgChan:
					_ = networkNode.BroadcastMessage(ctx, msg)
					return msg, nil
				}
			},
		),
		db,
		layer2.NewKeysContainer(),
		queue.NewQBProblemProcessor(
			queue.NewSettings(&queue.SSettings{
				FMessageConstructSettings: layer1.NewConstructSettings(&layer1.SConstructSettings{
					FSettings: layer1.NewSettings(&layer1.SSettings{
						FWorkSizeBits: tcWorkSize,
					}),
					FParallel: parallel,
				}),
				FQueuePoolCap: [2]uint64{tcQueueCap, tcQueueCap},
				FQueuePeriod:  time.Second,
				FConsumersCap: 1,
			}),
			func() layer2.IScheme {
				scheme, _ := hybrid.NewScheme(privKey, tcMsgSize)
				return scheme
			}(),
		),
	)
	go func() { _ = node.Run(ctx) }()
	return node, networkNode, privKey
}

func testRunNode(ctx context.Context, timeWait time.Duration, addr string, typeDB, numDB int, handlerF IHandlerF) (INode, network.INode, asymmetric.IPrivKey) {
	db, err := database.NewKVDatabase(fmt.Sprintf(tcPathDBTemplate, typeDB, numDB))
	if err != nil {
		panic(err)
	}
	return testRunNodeWithDB(ctx, timeWait, addr, db, handlerF)
}

func testFreeNodes(nodes []INode, typeDB int) {
	for _, node := range nodes {
		_ = node.GetKVDatabase().Close()
	}
	testDeleteDB(typeDB)
}

func testDeleteDB(typeDB int) {
	for i := 0; i < 5; i++ {
		_ = os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, typeDB, i))
	}
}

func (p *sNode) testNewNetworkMessage(pSett layer1.IConstructSettings, pMsgBytes []byte) layer1.IMessage {
	return layer1.NewMessage(
		pSett,
		pMsgBytes,
	)
}

func (p *sNode) testNewNetworkMessageWithAnotherMessageType(pSett layer1.IConstructSettings, pMsgBytes []byte) layer1.IMessage {
	return layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{
				FNetworkKey: pSett.GetSettings().GetNetworkKey() + "_",
			}),
		}),
		pMsgBytes,
	)
}

type tsDatabase struct{}

func (p *tsDatabase) Get([]byte) ([]byte, error) { return nil, database.ErrNotFound }
func (p *tsDatabase) Set([]byte, []byte) error   { return errors.New("some error") } //nolint:err113
func (p *tsDatabase) Del([]byte) error           { return nil }
func (p *tsDatabase) Close() error               { return nil }
