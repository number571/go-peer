package conn

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcHead = 12345
	tcBody = "hello, world!"
)

func TestConn(t *testing.T) {
	listener := testNewService(t)
	defer testFreeService(listener)

	conn, err := NewConn(
		NewSettings(&SSettings{
			FNetworkKey:    "_",
			FMessageSize:   testutils.TCMessageSize,
			FLimitVoidSize: 1, // not used
			FFetchTimeWait: 5 * time.Second,
		}),
		testutils.TgAddrs[17],
	)
	if err != nil {
		t.Error(err)
		return
	}

	pld, err := conn.FetchPayload(payload.NewPayload(tcHead, []byte(tcBody)))
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(pld.GetBody(), []byte(tcBody)) {
		t.Error("load payload not equal new payload")
		return
	}
}

func testNewService(t *testing.T) net.Listener {
	listener, err := net.Listen("tcp", testutils.TgAddrs[17])
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
					FNetworkKey:    "_",
					FMessageSize:   testutils.TCMessageSize,
					FLimitVoidSize: 1,
					FFetchTimeWait: 5 * time.Second,
				}),
				aconn,
			)
			pld := conn.ReadPayload()
			if pld == nil {
				break
			}

			ok := func() bool {
				defer conn.Close()
				return conn.WritePayload(pld) == nil
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
