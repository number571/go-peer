package network

import (
	"bytes"
	"net"
	"sync"
	"time"
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

func (conn *sConn) Socket() net.Conn {
	return conn.fSocket
}

func (conn *sConn) Request(msg IMessage) IMessage {
	var (
		chMsg    = make(chan IMessage)
		timeWait = conn.fSettings.GetTimeWait()
	)

	conn.Write(msg)
	go readMessage(conn, chMsg)

	select {
	case rmsg := <-chMsg:
		return rmsg
	case <-time.After(timeWait):
		return nil
	}
}

func (conn *sConn) Close() error {
	return conn.fSocket.Close()
}

func (conn *sConn) Write(msg IMessage) error {
	conn.fMutex.Lock()
	defer conn.fMutex.Unlock()

	msgBytes := msg.Bytes()
	ptr := len(msgBytes)

	for {
		n, err := conn.fSocket.Write(msgBytes[:ptr])
		if err != nil {
			return err
		}

		msgBytes = msgBytes[:n]
		ptr = ptr - n

		if ptr == 0 {
			break
		}
	}

	return nil
}

func (conn *sConn) Read() IMessage {
	chMsg := make(chan IMessage)
	go readMessage(conn, chMsg)
	return <-chMsg
}

func readMessage(conn *sConn, chMsg chan IMessage) {
	var msg IMessage
	defer func() {
		chMsg <- msg
	}()

	// bufLen = Size[u64] in bytes
	bufLen := make([]byte, cSizeUint)
	length, err := conn.fSocket.Read(bufLen)
	if err != nil {
		return
	}
	if length != cSizeUint {
		return
	}

	// mustLen = Size[u64] in uint64
	mustLen := newPackage(bufLen).BytesToSize()
	if mustLen > conn.fSettings.GetMessageSize() {
		return
	}

	msgRaw := make([]byte, 0, mustLen)
	for {
		buffer := make([]byte, mustLen)
		n, err := conn.fSocket.Read(buffer)
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
	msg = LoadMessage(bytes.Join(
		[][]byte{
			bufLen,
			msgRaw,
		},
		[]byte{},
	), []byte(conn.fSettings.GetNetworkKey()))
}
