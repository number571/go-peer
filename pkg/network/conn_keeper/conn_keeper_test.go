package conn_keeper

import (
	"net"
	"testing"
	"time"

	"github.com/number571/go-peer/internal/testutils"
	"github.com/number571/go-peer/pkg/network"
)

func TestConnKeeper(t *testing.T) {
	listener := testNewService(t)
	defer testFreeService(listener)

	node := network.NewNode(network.NewSettings(&network.SSettings{}))
	connKeeper := NewConnKeeper(
		NewSettings(&SSettings{
			FConnections: func() []string { return []string{testutils.TgAddrs[18]} },
			FDuration:    500 * time.Millisecond,
		}),
		node,
	)

	if err := connKeeper.Run(); err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second)
	if len(node.Connections()) != 1 {
		t.Error("lenght of connections != 1")
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