package conn

import (
	"bytes"
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcHead = 12345
	tcBody = "hello, world!"
)

func TestSettingsNetworkKey(t *testing.T) {
	t.Parallel()

	for i := 0; i < 5; i++ {
		testSettings(t, i)
	}
}

func TestSettings(t *testing.T) {
	t.Parallel()

	sett := NewSettings(&SSettings{
		FWorkSizeBits:          testutils.TCWorkSize,
		FLimitMessageSizeBytes: testutils.TCMessageSize,
		FWaitReadTimeout:       time.Hour,
		FDialTimeout:           time.Minute,
		FReadTimeout:           time.Minute,
		FWriteTimeout:          time.Minute,
	})

	networkKey := "hello, world!"
	sett.SetNetworkKey(networkKey)

	if sett.GetNetworkKey() != networkKey {
		t.Error("got invalid network key")
		return
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
			FWaitReadTimeout: time.Hour,
			FDialTimeout:     time.Minute,
			FReadTimeout:     time.Minute,
			FWriteTimeout:    time.Minute,
		})
	case 1:
		_ = NewSettings(&SSettings{
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		})
	case 2:
		_ = NewSettings(&SSettings{
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		})
	case 3:
		_ = NewSettings(&SSettings{
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
		})
	case 4:
		_ = NewSettings(&SSettings{
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		})
	}
}

func TestClosedConn(t *testing.T) {
	t.Parallel()

	listener := testNewService(t, testutils.TgAddrs[30], "")
	defer testFreeService(listener)

	conn, err := NewConn(
		NewSettings(&SSettings{
			FWorkSizeBits:          testutils.TCWorkSize,
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		}),
		testutils.TgAddrs[30],
	)
	if err != nil {
		t.Error(err)
		return
	}

	if err := conn.Close(); err != nil {
		t.Error(err)
		return
	}

	pld := payload.NewPayload(1, []byte("aaa"))
	msg := message.NewMessage(conn.GetSettings(), pld, 1, 0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := conn.WriteMessage(ctx, msg); err == nil {
		t.Error("success write payload to closed connection")
		return
	}

	readCh := make(chan struct{})
	go func() { <-readCh }()

	if _, err := conn.ReadMessage(ctx, readCh); err == nil {
		t.Error("success read payload from closed connection")
		return
	}

	sconn := conn.(*sConn)
	if err := sconn.sendBytes(ctx, []byte("hello, world!")); err == nil {
		t.Error("success send bytes to closed connection")
		return
	}

	if _, err := sconn.recvDataBytes(ctx, 128, time.Second); err == nil {
		t.Error("success recv data bytes from closed connection")
		return
	}
}

func TestInvalidConn(t *testing.T) {
	t.Parallel()

	_, err := NewConn(
		NewSettings(&SSettings{
			FWorkSizeBits:          testutils.TCWorkSize,
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		}),
		"INVALID_ADDRESS",
	)
	if err == nil {
		t.Error("success connect to invalid address")
		return
	}
}

func TestConnWithNetworkKey(t *testing.T) {
	t.Parallel()

	testConn(t, testutils.TgAddrs[17], "")
	// testConn(t, testutils.TgAddrs[17], "hello, world!")
}

func testConn(t *testing.T, pAddr, pNetworkKey string) {
	listener := testNewService(t, pAddr, pNetworkKey)
	defer testFreeService(listener)

	conn, err := NewConn(
		NewSettings(&SSettings{
			FWorkSizeBits:          testutils.TCWorkSize,
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		}),
		pAddr,
	)
	if err != nil {
		t.Error(err)
		return
	}

	socket := conn.GetSocket()
	remoteAddr := strings.ReplaceAll(pAddr, "localhost", "127.0.0.1")
	if socket.RemoteAddr().String() != remoteAddr {
		t.Error("got incorrect remote address")
		return
	}

	conn.GetSettings().SetNetworkKey(pNetworkKey)

	pld := payload.NewPayload(tcHead, []byte(tcBody))
	msg := message.NewMessage(conn.GetSettings(), pld, 1, 0)
	ctx := context.Background()
	if err := conn.WriteMessage(ctx, msg); err != nil {
		t.Error(err)
		return
	}

	readCh := make(chan struct{})
	go func() { <-readCh }()

	msgRecv, err := conn.ReadMessage(ctx, readCh)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(msgRecv.GetPayload().GetBody(), []byte(tcBody)) {
		t.Error("load payload not equal new payload")
		return
	}
}

func testNewService(t *testing.T, pAddr, pNetworkKey string) net.Listener {
	listener, err := net.Listen("tcp", pAddr)
	if err != nil {
		t.Error(err)
		return nil
	}

	go func() {
		for {
			aconn, err := listener.Accept()
			if err != nil {
				break
			}

			conn := LoadConn(
				NewSettings(&SSettings{
					FWorkSizeBits:          testutils.TCWorkSize,
					FLimitMessageSizeBytes: testutils.TCMessageSize,
					FWaitReadTimeout:       time.Hour,
					FDialTimeout:           time.Minute,
					FReadTimeout:           time.Minute,
					FWriteTimeout:          time.Minute,
				}),
				aconn,
			)

			conn.GetSettings().SetNetworkKey(pNetworkKey)

			readCh := make(chan struct{})
			go func() { <-readCh }()

			ctx := context.Background()

			msg, err := conn.ReadMessage(ctx, readCh)
			if err != nil {
				break
			}

			ok := func() bool {
				defer conn.Close()
				return conn.WriteMessage(ctx, msg) == nil
			}()

			if !ok {
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
