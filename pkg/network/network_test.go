package network

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcIter     = 100
	tcTimeWait = time.Minute
)

func TestSettings(t *testing.T) {
	for i := 0; i < 5; i++ {
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
			FMaxConnects:  testutils.TCMaxConnects,
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
			FCapacity:     testutils.TCCapacity,
			FReadTimeout:  tcTimeWait,
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
			FAddress:      "test",
			FCapacity:     testutils.TCCapacity,
			FMaxConnects:  testutils.TCMaxConnects,
			FWriteTimeout: tcTimeWait,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadDeadline: time.Hour,
				FReadDeadline:     time.Minute,
				FWriteDeadline:    time.Minute,
			}),
		})
	case 3:
		_ = NewSettings(&SSettings{
			FAddress:     "test",
			FCapacity:    testutils.TCCapacity,
			FMaxConnects: testutils.TCMaxConnects,
			FReadTimeout: tcTimeWait,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadDeadline: time.Hour,
				FReadDeadline:     time.Minute,
				FWriteDeadline:    time.Minute,
			}),
		})
	case 4:
		_ = NewSettings(&SSettings{
			FAddress:      "test",
			FCapacity:     testutils.TCCapacity,
			FMaxConnects:  testutils.TCMaxConnects,
			FReadTimeout:  tcTimeWait,
			FWriteTimeout: tcTimeWait,
		})
	}
}

func TestBroadcast(t *testing.T) {
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
	handleF := func(node INode, conn conn.IConn, reqBytes []byte) error {
		defer wg.Done()
		defer node.BroadcastPayload(payload.NewPayload(headHandle, reqBytes))

		tcMutex.Lock()
		defer tcMutex.Unlock()

		val := string(reqBytes)
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
	for i := 0; i < tcIter; i++ {
		go func(i int) {
			pld := payload.NewPayload(
				headHandle,
				[]byte(fmt.Sprintf(testutils.TcBodyTemplate, i)),
			)
			nodes[0].BroadcastPayload(pld)
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
	var (
		node1 = newTestNode("", 2).(*sNode)
		node2 = newTestNode(testutils.TgAddrs[27], 1)
		node3 = newTestNode(testutils.TgAddrs[28], testutils.TCMaxConnects)
	)

	if err := node2.Run(); err != nil {
		t.Error(err)
		return
	}
	defer node2.Stop()

	if err := node2.Run(); err == nil {
		t.Error("success second run node")
		return
	}

	if err := node3.Run(); err != nil {
		t.Error(err)
		return
	}
	defer node3.Stop()

	if err := node1.AddConnection("unknown_connection_address"); err == nil {
		t.Error("success add incorrect connection address")
		return
	}

	if err := node1.AddConnection(testutils.TgAddrs[27]); err != nil {
		t.Error(err)
		return
	}

	if err := node1.AddConnection(testutils.TgAddrs[27]); err == nil {
		t.Error("success add already exist connection")
		return
	}

	if err := node1.AddConnection(testutils.TgAddrs[28]); err != nil {
		t.Error(err)
		return
	}

	if err := node1.AddConnection(testutils.TgAddrs[28]); err == nil {
		t.Error("success add second connection with limit = 1")
		return
	}

	if err := node3.AddConnection(testutils.TgAddrs[27]); err != nil {
		t.Error(err)
		return
	}

	time.Sleep(200 * time.Millisecond)
	if len(node3.GetConnections()) != 1 {
		t.Error("has more than 1 connections (node2 should be auto disconnects by max conns param)")
		return
	}

	for _, c := range node1.GetConnections() {
		c.Close()
	}

	if err := node1.DelConnection(testutils.TgAddrs[27]); err == nil {
		t.Error("success delete already closed connection")
		return
	}
}

func TestHandleMessage(t *testing.T) {
	node := newTestNode("", testutils.TCMaxConnects).(*sNode)

	node.HandleFunc(1, nil)
	if ok := node.handleMessage(nil, payload.NewPayload(1, []byte{1})); ok {
		t.Error("success handle message with nil function")
		return
	}

	node.HandleFunc(1, func(i1 INode, i2 conn.IConn, b []byte) error {
		return errors.New("some error")
	})
	if ok := node.handleMessage(nil, payload.NewPayload(1, []byte{2})); ok {
		t.Error("success handle message with got error from function")
		return
	}

	node.HandleFunc(1, func(i1 INode, i2 conn.IConn, b []byte) error {
		return nil
	})
	if ok := node.handleMessage(nil, payload.NewPayload(1, []byte{3})); !ok {
		t.Error("failed handle message with correct function")
		return
	}
}

func TestNodeSettings(t *testing.T) {
	gotSett := newTestNode("", testutils.TCMaxConnects).GetSettings()
	if gotSett.GetMaxConnects() != testutils.TCMaxConnects {
		t.Error("invalid setting's value")
	}
}

func testNodes() ([5]INode, map[INode]map[string]bool, error) {
	nodes := [5]INode{}
	addrs := [5]string{"", "", testutils.TgAddrs[0], "", testutils.TgAddrs[1]}

	for i := 0; i < 5; i++ {
		nodes[i] = newTestNode(addrs[i], testutils.TCMaxConnects)
	}

	if err := nodes[2].Run(); err != nil {
		return nodes, nil, err
	}
	if err := nodes[4].Run(); err != nil {
		return nodes, nil, err
	}

	time.Sleep(500 * time.Millisecond)

	nodes[0].AddConnection(testutils.TgAddrs[0])
	nodes[1].AddConnection(testutils.TgAddrs[1])

	nodes[3].AddConnection(testutils.TgAddrs[0])
	nodes[3].AddConnection(testutils.TgAddrs[1])

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
	sett := NewSettings(&SSettings{
		FAddress:      pAddr,
		FCapacity:     testutils.TCCapacity,
		FMaxConnects:  pMaxConns,
		FReadTimeout:  time.Minute,
		FWriteTimeout: time.Minute,
		FConnSettings: conn.NewSettings(&conn.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadDeadline: time.Hour,
			FReadDeadline:     time.Minute,
			FWriteDeadline:    time.Minute,
		}),
	})
	return NewNode(sett)
}

func testFreeNodes(nodes []INode) {
	for _, node := range nodes {
		node.Stop()
	}
}
