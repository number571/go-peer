// nolint: goerr113
package connkeeper

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/cache"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SConnKeeperError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestSettings(t *testing.T) {
	t.Parallel()

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
			FConnections: func() []string { return []string{testutils.TgAddrs[7]} },
		})
	}
}

func TestConnKeeperSettings(t *testing.T) {
	t.Parallel()

	duration := time.Second / 2
	connKeeper := newTestConnKeeper(duration)

	if connKeeper.GetSettings().GetDuration() != duration {
		t.Error("got invalid settings param")
		return
	}
}

func TestConnKeeper(t *testing.T) {
	t.Parallel()

	listener := testNewService(t)
	defer testFreeService(listener)

	connKeeper := newTestConnKeeper(50 * time.Millisecond)
	if node := connKeeper.GetNetworkNode(); node == nil {
		t.Error("network node is nil")
		return
	}

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer func() {
		cancel1()
		time.Sleep(100 * time.Millisecond)
	}()

	go func() {
		if err := connKeeper.Run(ctx1); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	go func() {
		err1 := testutils.TryN(50, 20*time.Millisecond, func() error {
			if err := connKeeper.Run(ctx2); err == nil {
				return errors.New("error is nil with already running connKeeper")
			}
			return nil
		})
		if err1 != nil {
			t.Error(err1)
			return
		}
	}()

	err1 := testutils.TryN(50, 20*time.Millisecond, func() error {
		if len(connKeeper.GetNetworkNode().GetConnections()) != 1 {
			return errors.New("length of connections != 1")
		}
		return nil
	})
	if err1 != nil {
		t.Error(err1)
		return
	}
}

func newTestConnKeeper(pDuration time.Duration) IConnKeeper {
	return NewConnKeeper(
		NewSettings(&SSettings{
			FConnections: func() []string { return []string{testutils.TgAddrs[7]} },
			FDuration:    pDuration,
		}),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FMaxConnects:  16,
				FReadTimeout:  time.Minute,
				FWriteTimeout: time.Minute,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FMessageSettings: layer1.NewSettings(&layer1.SSettings{
						FWorkSizeBits: 10,
					}),
					FLimitMessageSizeBytes: (8 << 10),
					FWaitReadTimeout:       time.Hour,
					FDialTimeout:           time.Minute,
					FReadTimeout:           time.Minute,
					FWriteTimeout:          time.Minute,
				}),
			}),
			cache.NewLRUCache(1024),
		),
	)
}

func testNewService(t *testing.T) net.Listener {
	listener, err := net.Listen("tcp", testutils.TgAddrs[7])
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
