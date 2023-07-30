package conn_keeper

import (
	"net"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestConnKeeper(t *testing.T) {
	listener := testNewService(t)
	defer testFreeService(listener)

	node := network.NewNode(network.NewSettings(&network.SSettings{
		FCapacity:     testutils.TCCapacity,
		FMaxConnects:  testutils.TCMaxConnects,
		FWriteTimeout: time.Minute,
		FConnSettings: conn.NewSettings(&conn.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadDeadline: time.Hour,
			FReadDeadline:     time.Minute,
			FWriteDeadline:    time.Minute,
			FFetchTimeWait:    1, // not used
		}),
	}))
	connKeeper := NewConnKeeper(
		NewSettings(&SSettings{
			FConnections: func() []string { return []string{testutils.TgAddrs[18]} },
			FDuration:    500 * time.Millisecond,
		}),
		node,
	)

	if node := connKeeper.GetNetworkNode(); node == nil {
		t.Error("network node is nil")
		return
	}

	if err := connKeeper.Run(); err != nil {
		t.Error(err)
		return
	}

	if err := connKeeper.Run(); err == nil {
		t.Error("error is nil with already running connKeeper")
		return
	}

	time.Sleep(time.Second)
	if len(node.GetConnections()) != 1 {
		t.Error("length of connections != 1")
		return
	}

	if err := connKeeper.Stop(); err != nil {
		t.Error(err)
		return
	}

	if err := connKeeper.Stop(); err == nil {
		t.Error("error is nil with already closed connKeeper")
		return
	}
}

func testNewService(t *testing.T) net.Listener {
	listener, err := net.Listen("tcp", testutils.TgAddrs[18])
	if err != nil {
		t.Error(err)
		return nil
	}

	go func() {
		for {
			_, err := listener.Accept()
			if err != nil {
				break
			}
		}
	}()

	return listener
}

func testFreeService(listener net.Listener) {
	if listener == nil {
		return
	}
	listener.Close()
}
