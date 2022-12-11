package conn

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/number571/go-peer/modules/encoding"
	"github.com/number571/go-peer/modules/network/message"
	"github.com/number571/go-peer/modules/payload"
)

var (
	_ IConn = &sConn{}
)

type sConn struct {
	fMutex    sync.Mutex
	fSocket   net.Conn
	fSettings ISettings
}

func NewConn(sett ISettings, address string) (IConn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return LoadConn(sett, conn), nil
}

func LoadConn(sett ISettings, conn net.Conn) IConn {
	return &sConn{
		fSettings: sett,
		fSocket:   conn,
	}
}

func (conn *sConn) Settings() ISettings {
	return conn.fSettings
}

func (conn *sConn) Socket() net.Conn {
	return conn.fSocket
}

func (conn *sConn) Request(pld payload.IPayload) (payload.IPayload, error) {
	var (
		chPld    = make(chan payload.IPayload)
		timeWait = conn.fSettings.GetTimeWait()
	)

	if err := conn.Write(pld); err != nil {
		return nil, err
	}
	go readPayload(conn, chPld)

	select {
	case rpld := <-chPld:
		if rpld == nil {
			return nil, fmt.Errorf("failed: read payload")
		}
		return rpld, nil
	case <-time.After(timeWait):
		return nil, fmt.Errorf("failed: time out")
	}
}

func (conn *sConn) Close() error {
	return conn.Socket().Close()
}

func (conn *sConn) Write(pld payload.IPayload) error {
	conn.fMutex.Lock()
	defer conn.fMutex.Unlock()

	var (
		msgBytes  = message.NewMessage(pld, []byte(conn.fSettings.GetNetworkKey())).Bytes()
		packBytes = payload.NewPayload(uint64(len(msgBytes)), msgBytes).Bytes()
		packPtr   = len(packBytes)
	)

	for {
		n, err := conn.Socket().Write(packBytes[:packPtr])
		if err != nil {
			return err
		}

		packPtr = packPtr - n
		packBytes = packBytes[:packPtr]

		if packPtr == 0 {
			break
		}
	}

	return nil
}

func (conn *sConn) Read() payload.IPayload {
	chPld := make(chan payload.IPayload)
	go readPayload(conn, chPld)
	return <-chPld
}

func readPayload(conn *sConn, chPld chan payload.IPayload) {
	var pld payload.IPayload
	defer func() {
		chPld <- pld
	}()

	// bufLen = Size[u64] in bytes
	bufLen := make([]byte, encoding.CSizeUint64)
	length, err := conn.Socket().Read(bufLen)
	if err != nil {
		return
	}
	if length != encoding.CSizeUint64 {
		return
	}

	// mustLen = Size[u64] in uint64
	arrLen := [encoding.CSizeUint64]byte{}
	copy(arrLen[:], bufLen)

	mustLen := encoding.BytesToUint64(arrLen)
	if mustLen > conn.fSettings.GetMessageSize() {
		return
	}

	msgRaw := make([]byte, 0, mustLen)
	for {
		buffer := make([]byte, mustLen)
		n, err := conn.Socket().Read(buffer)
		if err != nil {
			return
		}

		msgRaw = bytes.Join(
			[][]byte{
				msgRaw,
				buffer[:n],
			},
			[]byte{},
		)

		mustLen -= uint64(n)
		if mustLen == 0 {
			break
		}
	}

	// try unpack message from bytes
	msg := message.LoadMessage(
		msgRaw,
		[]byte(conn.fSettings.GetNetworkKey()),
	)
	if msg == nil {
		return
	}

	pld = msg.Payload()
}
