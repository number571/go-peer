package conn_keeper

import (
	"net"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestSettings(t *testing.T) {
	for i := 0; i < 2; i++ {
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
			FDuration: 500 * time.Millisecond,
		})
	case 1:
		_ = NewSettings(&SSettings{
			FConnections: func() []string { return []string{testutils.TgAddrs[18]} },
		})
	}
}

func TestConnKeeperSettings(t *testing.T) {
	duration := time.Second / 2
	connKeeper := newTestConnKeeper(duration)

	if connKeeper.GetSettings().GetDuration() != duration {
		t.Error("got invalid settings param")
		return
	}
}

func TestConnKeeper(t *testing.T) {
	listener := testNewService(t)
	defer testFreeService(listener)

	connKeeper := newTestConnKeeper(time.Second / 2)

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
	if len(connKeeper.GetNetworkNode().GetConnections()) != 1 {
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

func newTestConnKeeper(pDuration time.Duration) IConnKeeper {
	return NewConnKeeper(
		NewSettings(&SSettings{
			FConnections: func() []string { return []string{testutils.TgAddrs[18]} },
			FDuration:    pDuration,
		}),
		network.NewNode(network.NewSettings(&network.SSettings{
			FCapacity:     testutils.TCCapacity,
			FMaxConnects:  testutils.TCMaxConnects,
			FReadTimeout:  time.Minute,
			FWriteTimeout: time.Minute,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadDeadline: time.Hour,
				FReadDeadline:     time.Minute,
				FWriteDeadline:    time.Minute,
			}),
		})),
	)
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
