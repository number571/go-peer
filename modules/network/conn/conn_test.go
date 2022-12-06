package conn

import (
	"bytes"
	"net"
	"testing"

	"github.com/number571/go-peer/modules/payload"
	"github.com/number571/go-peer/settings/testutils"
)

const (
	tcHead = 12345
	tcBody = "hello, world!"
)

func TestConn(t *testing.T) {
	listener := testNewService(t)
	defer testFreeService(listener)

	conn, err := NewConn(NewSettings(&SSettings{}), testutils.TgAddrs[17])
	if err != nil {
		t.Error(err)
		return
	}

	pld, err := conn.Request(payload.NewPayload(tcHead, []byte(tcBody)))
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(pld.Body(), []byte(tcBody)) {
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

			conn := LoadConn(NewSettings(&SSettings{}), aconn)
			pld := conn.Read()

			ok := func() bool {
				defer conn.Close()
				return conn.Write(pld) == nil
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
