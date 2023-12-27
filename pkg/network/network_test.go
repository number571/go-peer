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
	"github.com/number571/go-peer/pkg/queue_set"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcIter     = 100
	tcTimeWait = time.Minute
)

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
				FMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadDeadline: time.Hour,
				FReadDeadline:     time.Minute,
				FWriteDeadline:    time.Minute,
			}),
		})
	case 1:
		_ = NewSettings(&SSettings{
			FAddress:      "test",
			FMaxConnects:  testutils.TCMaxConnects,
			FWriteTimeout: tcTimeWait,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadDeadline: time.Hour,
				FReadDeadline:     time.Minute,
				FWriteDeadline:    time.Minute,
			}),
		})
	case 2:
		_ = NewSettings(&SSettings{
			FAddress:     "test",
			FMaxConnects: testutils.TCMaxConnects,
			FReadTimeout: tcTimeWait,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadDeadline: time.Hour,
				FReadDeadline:     time.Minute,
				FWriteDeadline:    time.Minute,
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

	headHandle := uint64(testutils.TcHead)
	handleF := func(pCtx context.Context, node INode, conn conn.IConn, pMsg message.IMessage) error {
		defer wg.Done()
		defer node.BroadcastMessage(pCtx, pMsg)

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
			pld := payload.NewPayload(
				headHandle,
				[]byte(fmt.Sprintf(testutils.TcBodyTemplate, i)),
			)
			sett := nodes[0].GetSettings().GetConnSettings()
			nodes[0].BroadcastMessage(ctx, message.NewMessage(sett, pld))
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

// func newListener(t *testing.T, addr string) net.Listener {
// 	listener, err := net.Listen("tcp", addr)
// 	if err != nil {
// 		t.Error(err)
// 		return nil
// 	}
// 	go func() {
// 		for {
// 			conn, err := listener.Accept()
// 			if err != nil {
// 				return
// 			}
// 			_ = conn.Close()
// 		}
// 	}()
// 	return listener
// }

// func TestClosedConnection(t *testing.T) {
// 	var (
// 		node1    = newTestNode("", 1, time.Minute).(*sNode)
// 		listener = newListener(t, testutils.TgAddrs[37])
// 	)
// 	defer testFreeNodes([]INode{node1})

// 	defer func() {
// 		if listener == nil {
// 			return
// 		}
// 		listener.Close()
// 	}()

// 	if err := node1.AddConnection(testutils.TgAddrs[37]); err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	headHandle := uint64(testutils.TcHead)
// 	reqBytes := []byte("hello, world!")

// 	pld := payload.NewPayload(headHandle, reqBytes)
// 	if err := node1.BroadcastPayload(pld); err == nil {
// 		t.Error("success broadcast payload with non listening server")
// 		return
// 	}
// }

func TestNodeConnection(t *testing.T) {
	t.Parallel()

	var (
		node1 = newTestNode("", 2, time.Minute).(*sNode)
		node2 = newTestNode(testutils.TgAddrs[27], 1, time.Minute)
		node3 = newTestNode(testutils.TgAddrs[28], testutils.TCMaxConnects, time.Minute)
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

	for _, c := range node1.GetConnections() {
		c.Close()
	}

	if err := node1.DelConnection(testutils.TgAddrs[27]); err == nil {
		t.Error("success delete already closed connection")
		return
	}

	if err := node2.Close(); err != nil {
		t.Error(err)
		return
	}
}

func TestHandleMessage(t *testing.T) {
	t.Parallel()

	node := newTestNode("", testutils.TCMaxConnects, time.Minute).(*sNode)
	defer testFreeNodes([]INode{node})

	sett := node.GetSettings().GetConnSettings()
	ctx := context.Background()

	node.HandleFunc(1, nil)
	msg1 := message.NewMessage(sett, payload.NewPayload(1, []byte{1}))
	if ok := node.handleMessage(ctx, nil, msg1); ok {
		t.Error("success handle message with nil function")
		return
	}

	node.HandleFunc(1, func(ctx context.Context, i1 INode, i2 conn.IConn, b message.IMessage) error {
		return errors.New("some error")
	})
	msg2 := message.NewMessage(sett, payload.NewPayload(1, []byte{2}))
	if ok := node.handleMessage(ctx, nil, msg2); ok {
		t.Error("success handle message with got error from function")
		return
	}

	node.HandleFunc(1, func(ctx context.Context, i1 INode, i2 conn.IConn, b message.IMessage) error {
		return nil
	})
	msg3 := message.NewMessage(sett, payload.NewPayload(1, []byte{3}))
	if ok := node.handleMessage(ctx, nil, msg3); !ok {
		t.Error("failed handle message with correct function")
		return
	}
}

func TestNodeSettings(t *testing.T) {
	t.Parallel()

	gotSett := newTestNode("", testutils.TCMaxConnects, time.Minute).GetSettings()
	if gotSett.GetMaxConnects() != testutils.TCMaxConnects {
		t.Error("invalid setting's value")
	}
}

func testNodes() ([5]INode, map[INode]map[string]bool, error) {
	nodes := [5]INode{}
	addrs := [5]string{"", "", testutils.TgAddrs[0], "", testutils.TgAddrs[1]}

	for i := 0; i < 5; i++ {
		nodes[i] = newTestNode(addrs[i], testutils.TCMaxConnects, time.Minute)
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

	nodes[3].AddConnection(ctx, testutils.TgAddrs[0])
	nodes[3].AddConnection(ctx, testutils.TgAddrs[1])

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

func newTestNode(pAddr string, pMaxConns uint64, timeout time.Duration) INode {
	return NewNode(
		NewSettings(&SSettings{
			FAddress:      pAddr,
			FMaxConnects:  pMaxConns,
			FReadTimeout:  timeout,
			FWriteTimeout: timeout,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FWorkSizeBits:     testutils.TCWorkSize,
				FMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadDeadline: time.Hour,
				FReadDeadline:     timeout,
				FWriteDeadline:    timeout,
			}),
		}),
		queue_set.NewQueueSet(
			queue_set.NewSettings(&queue_set.SSettings{
				FCapacity: testutils.TCCapacity,
			}),
		),
	)
}

func testFreeNodes(nodes []INode) {
	for _, node := range nodes {
		node.Close()
	}
}
