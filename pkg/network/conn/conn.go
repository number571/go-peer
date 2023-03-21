package conn

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
)

var (
	_ IConn = &sConn{}
)

type sConn struct {
	fMutex    sync.Mutex
	fSocket   net.Conn
	fSettings ISettings
}

func NewConn(pSett ISettings, pAddr string) (IConn, error) {
	conn, err := net.Dial("tcp", pAddr)
	if err != nil {
		return nil, err
	}
	return LoadConn(pSett, conn), nil
}

func LoadConn(pSett ISettings, pConn net.Conn) IConn {
	return &sConn{
		fSettings: pSett,
		fSocket:   pConn,
	}
}

func (p *sConn) GetSettings() ISettings {
	return p.fSettings
}

func (p *sConn) GetSocket() net.Conn {
	return p.fSocket
}

func (p *sConn) FetchPayload(pPld payload.IPayload) (payload.IPayload, error) {
	var (
		chPld    = make(chan payload.IPayload)
		timeWait = p.fSettings.GetTimeWait()
	)

	if err := p.WritePayload(pPld); err != nil {
		return nil, err
	}
	go readPayload(p, chPld)

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

func (p *sConn) Close() error {
	return p.GetSocket().Close()
}

func (p *sConn) WritePayload(pPld payload.IPayload) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	var (
		msgBytes  = message.NewMessage(pPld, []byte(p.fSettings.GetNetworkKey())).GetBytes()
		packBytes = payload.NewPayload(uint64(len(msgBytes)), msgBytes).ToBytes()
		packPtr   = len(packBytes)
	)

	for {
		n, err := p.GetSocket().Write(packBytes[:packPtr])
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

func (p *sConn) ReadPayload() payload.IPayload {
	chPld := make(chan payload.IPayload)
	go readPayload(p, chPld)
	return <-chPld
}

func readPayload(pConn *sConn, pChPld chan payload.IPayload) {
	var pld payload.IPayload
	defer func() {
		pChPld <- pld
	}()

	// bufLen = Size[u64] in bytes
	bufLen := make([]byte, encoding.CSizeUint64)
	length, err := pConn.GetSocket().Read(bufLen)
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
	if mustLen > pConn.fSettings.GetMessageSize() {
		return
	}

	msgRaw := make([]byte, 0, mustLen)
	for {
		buffer := make([]byte, mustLen)
		n, err := pConn.GetSocket().Read(buffer)
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
		[]byte(pConn.fSettings.GetNetworkKey()),
	)
	if msg == nil {
		return
	}

	pld = msg.GetPayload()
}
