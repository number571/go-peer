// nolint: goerr113
package network

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache/lru"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcIter     = 100
	tcTimeWait = time.Minute
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SNetworkError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestSettings(t *testing.T) {
	t.Parallel()

	for i := 0; i < 4; i++ {
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
			FAddress:      "test",
			FReadTimeout:  tcTimeWait,
			FWriteTimeout: tcTimeWait,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FLimitMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadTimeout:       time.Hour,
				FDialTimeout:           time.Minute,
				FReadTimeout:           time.Minute,
				FWriteTimeout:          time.Minute,
			}),
		})
	case 1:
		_ = NewSettings(&SSettings{
			FAddress:      "test",
			FMaxConnects:  testutils.TCMaxConnects,
			FWriteTimeout: tcTimeWait,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FLimitMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadTimeout:       time.Hour,
				FDialTimeout:           time.Minute,
				FReadTimeout:           time.Minute,
				FWriteTimeout:          time.Minute,
			}),
		})
	case 2:
		_ = NewSettings(&SSettings{
			FAddress:     "test",
			FMaxConnects: testutils.TCMaxConnects,
			FReadTimeout: tcTimeWait,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FLimitMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadTimeout:       time.Hour,
				FDialTimeout:           time.Minute,
				FReadTimeout:           time.Minute,
				FWriteTimeout:          time.Minute,
			}),
		})
	case 3:
		_ = NewSettings(&SSettings{
			FAddress:      "test",
			FMaxConnects:  testutils.TCMaxConnects,
			FReadTimeout:  tcTimeWait,
			FWriteTimeout: tcTimeWait,
		})
	}
}

func TestBroadcast(t *testing.T) {
	t.Parallel()

	nodes, mapp, err := testNodes()
	if err != nil {
		t.Error(err)
		return
	}
	defer testFreeNodes(nodes[:])

	// four receivers, sender not receive his messages
	tcMutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(4 * tcIter)

	headHandle := testutils.TcHead
	handleF := func(pCtx context.Context, node INode, _ conn.IConn, pMsg message.IMessage) error {
		defer func() {
			_ = node.BroadcastMessage(pCtx, pMsg)
			wg.Done()
		}()

		tcMutex.Lock()
		defer tcMutex.Unlock()

		val := string(pMsg.GetPayload().GetBody())
		flag, ok := mapp[node][val]
		if !ok {
			err := fmt.Errorf("incoming value '%s' undefined", val)
			t.Error(err)
			return err
		}

		if flag {
			err := fmt.Errorf("incoming value '%s' already exists", val)
			t.Error(err)
			return err
		}

		mapp[node][val] = true
		return nil
	}

	for _, node := range nodes {
		node.HandleFunc(headHandle, handleF)
	}

	// nodes[0] -> nodes[1:]
	ctx := context.Background()
	for i := 0; i < tcIter; i++ {
		go func(i int) {
			pld := payload.NewPayload32(
				headHandle,
				[]byte(fmt.Sprintf(testutils.TcBodyTemplate, i)),
			)
			sett := message.NewSettings(&message.SSettings{
				FNetworkKey:   nodes[0].GetVSettings().GetNetworkKey(),
				FWorkSizeBits: nodes[0].GetSettings().GetConnSettings().GetWorkSizeBits(),
			})
			_ = nodes[0].BroadcastMessage(ctx, message.NewMessage(sett, pld))
		}(i)
	}

	ch := make(chan struct{})
	go func() {
		wg.Wait()
		ch <- struct{}{}
	}()

	select {
	case <-ch:
	case <-time.After(tcTimeWait):
		t.Error("limit of waiting time for group")
		return
	}

	for _, node := range nodes {
		// pass sender
		if node == nodes[0] {
			continue
		}
		for i := 0; i < tcIter; i++ {
			val := fmt.Sprintf(testutils.TcBodyTemplate, i)
			flag, ok := mapp[node][val]
			if !ok {
				t.Errorf("result value '%s' undefined", val)
				continue
			}
			if !flag {
				t.Errorf("result value '%s' not exists", val)
				continue
			}
		}
	}
}

func TestNodeConnection(t *testing.T) {
	t.Parallel()

	var (
		node1 = newTestNode("", 2).(*sNode)
		node2 = newTestNode(testutils.TgAddrs[27], 1)
		node3 = newTestNode(testutils.TgAddrs[28], testutils.TCMaxConnects)
	)
	defer testFreeNodes([]INode{node1, node2, node3})

	ctx := context.Background()
	go func() {
		if err := node2.Listen(ctx); err != nil && !errors.Is(err, net.ErrClosed) {
			t.Error(err)
			return
		}
	}()
	defer node2.Close()

	go func() {
		if err := node3.Listen(ctx); err != nil && !errors.Is(err, net.ErrClosed) {
			t.Error(err)
			return
		}
	}()
	defer node3.Close()

	time.Sleep(200 * time.Millisecond)
	go func() {
		if err := node2.Listen(ctx); err == nil {
			t.Error("success second run node")
			return
		}
	}()

	err1 := testutils.TryN(50, 10*time.Millisecond, func() error {
		if err := node1.AddConnection(ctx, "unknown_connection_address"); err == nil {
			return errors.New("success add incorrect connection address")
		}
		return nil
	})
	if err1 != nil {
		t.Error(err1)
		return
	}

	if err := node1.AddConnection(ctx, testutils.TgAddrs[27]); err != nil {
		t.Error(err)
		return
	}

	if err := node1.AddConnection(ctx, testutils.TgAddrs[27]); err == nil {
		t.Error("success add already exist connection")
		return
	}

	if err := node1.AddConnection(ctx, testutils.TgAddrs[28]); err != nil {
		t.Error(err)
		return
	}

	if err := node1.AddConnection(ctx, testutils.TgAddrs[28]); err == nil {
		t.Error("success add second connection with limit = 1")
		return
	}

	if err := node3.AddConnection(ctx, testutils.TgAddrs[27]); err != nil {
		t.Error(err)
		return
	}

	err2 := testutils.TryN(50, 10*time.Millisecond, func() error {
		if len(node3.GetConnections()) != 1 {
			return errors.New("has more than 1 connections (node2 should be auto disconnects by max conns param)")
		}
		return nil
	})
	if err2 != nil {
		t.Error(err2)
		return
	}

	node1.SetVSettings(conn.NewVSettings(&conn.SVSettings{
		FNetworkKey: "set_another_network_key",
	}))

	if err := node1.DelConnection(testutils.TgAddrs[27]); err == nil {
		t.Error("success delete already closed connection")
		return
	}

	if len(node1.GetConnections()) != 0 {
		t.Error("set message settings should close all connections")
		return
	}

	if err := node2.Close(); err != nil {
		t.Error(err)
		return
	}
}

func TestHandleMessage(t *testing.T) {
	t.Parallel()

	node := newTestNode("", testutils.TCMaxConnects).(*sNode)
	defer testFreeNodes([]INode{node})

	newNetworkKey := "handle_message_network_key"
	node.SetVSettings(conn.NewVSettings(&conn.SVSettings{
		FNetworkKey: newNetworkKey,
	}))

	ctx := context.Background()
	vsett := node.GetVSettings()

	if vsett.GetNetworkKey() != newNetworkKey {
		t.Error("incorrect set variable settings")
		return
	}

	sett := message.NewSettings(&message.SSettings{
		FNetworkKey:   vsett.GetNetworkKey(),
		FWorkSizeBits: node.GetSettings().GetConnSettings().GetWorkSizeBits(),
	})

	node.HandleFunc(1, nil)
	msg1 := message.NewMessage(sett, payload.NewPayload32(1, []byte{1}))
	if ok := node.handleMessage(ctx, nil, msg1); ok {
		t.Error("success handle message with nil function")
		return
	}

	node.HandleFunc(1, func(_ context.Context, _ INode, _ conn.IConn, _ message.IMessage) error {
		return errors.New("some error")
	})
	msg2 := message.NewMessage(sett, payload.NewPayload32(1, []byte{2}))
	if ok := node.handleMessage(ctx, nil, msg2); ok {
		t.Error("success handle message with got error from function")
		return
	}

	node.HandleFunc(1, func(_ context.Context, _ INode, _ conn.IConn, _ message.IMessage) error {
		return nil
	})
	msg3 := message.NewMessage(sett, payload.NewPayload32(1, []byte{3}))
	if ok := node.handleMessage(ctx, nil, msg3); !ok {
		t.Error("failed handle message with correct function")
		return
	}
}

func TestNodeSettings(t *testing.T) {
	t.Parallel()

	gotSett := newTestNode("", testutils.TCMaxConnects).GetSettings()
	if gotSett.GetMaxConnects() != testutils.TCMaxConnects {
		t.Error("invalid setting's value")
	}
}

func TestContextCancel(t *testing.T) {
	t.Parallel()

	node1 := newTestNode(testutils.TgAddrs[16], testutils.TCMaxConnects)
	node2 := newTestNode("", testutils.TCMaxConnects)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { _ = node1.Listen(ctx) }()

	err1 := testutils.TryN(50, 10*time.Millisecond, func() error {
		return node2.AddConnection(ctx, testutils.TgAddrs[16])
	})
	if err1 != nil {
		t.Error(err1)
		return
	}

	headHandle := testutils.TcHead
	sett := message.NewSettings(&message.SSettings{
		FNetworkKey:   node2.GetVSettings().GetNetworkKey(),
		FWorkSizeBits: node2.GetSettings().GetConnSettings().GetWorkSizeBits(),
	})

	go func() {
		for i := 0; i < 1000; i++ {
			pld := payload.NewPayload32(
				headHandle,
				[]byte(fmt.Sprintf(testutils.TcBodyTemplate, i)),
			)
			if err := node2.BroadcastMessage(ctx, message.NewMessage(sett, pld)); err != nil {
				return
			}
		}
		t.Error("success all broadcast messages with canceled context")
	}()

	cancel()
}

func testNodes() ([5]INode, map[INode]map[string]bool, error) {
	nodes := [5]INode{}
	addrs := [5]string{"", "", testutils.TgAddrs[0], "", testutils.TgAddrs[1]}

	for i := 0; i < 5; i++ {
		nodes[i] = newTestNode(addrs[i], testutils.TCMaxConnects)
	}

	ctx := context.Background()

	go func() { _ = nodes[2].Listen(ctx) }()
	go func() { _ = nodes[4].Listen(ctx) }()

	err1 := testutils.TryN(50, 10*time.Millisecond, func() error {
		return nodes[0].AddConnection(ctx, testutils.TgAddrs[0])
	})
	if err1 != nil {
		return [5]INode{}, nil, err1
	}
	err2 := testutils.TryN(50, 10*time.Millisecond, func() error {
		return nodes[1].AddConnection(ctx, testutils.TgAddrs[1])
	})
	if err2 != nil {
		return [5]INode{}, nil, err2
	}

	_ = nodes[3].AddConnection(ctx, testutils.TgAddrs[0])
	_ = nodes[3].AddConnection(ctx, testutils.TgAddrs[1])

	mapp := make(map[INode]map[string]bool)
	for _, node := range nodes {
		// pass sender
		if node == nodes[0] {
			continue
		}
		mapp[node] = make(map[string]bool)
		for i := 0; i < tcIter; i++ {
			mapp[node][fmt.Sprintf(testutils.TcBodyTemplate, i)] = false
		}
	}

	return nodes, mapp, nil
}

func newTestNode(pAddr string, pMaxConns uint64) INode {
	timeout := time.Minute
	return NewNode(
		NewSettings(&SSettings{
			FAddress:      pAddr,
			FMaxConnects:  pMaxConns,
			FReadTimeout:  timeout,
			FWriteTimeout: timeout,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FWorkSizeBits:          testutils.TCWorkSize,
				FLimitMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadTimeout:       time.Hour,
				FDialTimeout:           time.Minute,
				FReadTimeout:           timeout,
				FWriteTimeout:          timeout,
			}),
		}),
		conn.NewVSettings(&conn.SVSettings{}),
		lru.NewLRUCache(testutils.TCCapacity),
	)
}

func testFreeNodes(nodes []INode) {
	for _, node := range nodes {
		node.Close()
	}
}
