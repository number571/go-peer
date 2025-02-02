package conn

import (
	"bytes"
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/payload/joiner"
)

var (
	_ IConn = &sConn{}
)

type sConn struct {
	fMutex    sync.RWMutex
	fSocket   net.Conn
	fSettings ISettings
}

func Connect(pCtx context.Context, pSett ISettings, pAddr string) (IConn, error) {
	dialer := &net.Dialer{Timeout: pSett.GetDialTimeout()}
	conn, err := dialer.DialContext(pCtx, "tcp", pAddr)
	if err != nil {
		return nil, errors.Join(ErrCreateConnection, err)
	}
	return LoadConn(pSett, conn), nil
}

func LoadConn(pSett ISettings, pConn net.Conn) IConn {
	return &sConn{
		fSocket:   pConn,
		fSettings: pSett,
	}
}

func (p *sConn) GetSettings() ISettings {
	return p.fSettings
}

func (p *sConn) GetSocket() net.Conn {
	return p.fSocket
}

func (p *sConn) Close() error {
	return p.fSocket.Close()
}

func (p *sConn) WriteMessage(pCtx context.Context, pMsg layer1.IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	bytesJoiner := joiner.NewBytesJoiner32([][]byte{pMsg.ToBytes()})
	if err := p.sendBytes(pCtx, bytesJoiner); err != nil {
		return errors.Join(ErrSendPayloadBytes, err)
	}

	return nil
}

func (p *sConn) ReadMessage(pCtx context.Context, pChRead chan<- struct{}) (layer1.IMessage, error) {
	// large wait read deadline => the connection has not sent anything yet
	msgSize, err := p.recvHeadBytes(pCtx, pChRead, p.fSettings.GetWaitReadTimeout())
	if err != nil {
		return nil, errors.Join(ErrReadHeaderBytes, err)
	}

	dataBytes, err := p.recvDataBytes(pCtx, msgSize, p.fSettings.GetReadTimeout())
	if err != nil {
		return nil, errors.Join(ErrReadBodyBytes, err)
	}

	// try unpack message from bytes
	msg, err := layer1.LoadMessage(p.fSettings.GetMessageSettings(), dataBytes)
	if err != nil {
		return nil, errors.Join(ErrInvalidMessageBytes, err)
	}

	return msg, nil
}

func (p *sConn) sendBytes(pCtx context.Context, pBytes []byte) error {
	bytesPtr := uint64(len(pBytes))
	for bytesPtr != 0 {
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		default:
			if err := p.fSocket.SetWriteDeadline(time.Now().Add(p.fSettings.GetWriteTimeout())); err != nil {
				return errors.Join(ErrSetWriteDeadline, err)
			}

			n, err := p.fSocket.Write(pBytes[:bytesPtr])
			if err != nil {
				return errors.Join(ErrWriteToSocket, err)
			}

			bytesPtr -= uint64(n)
			pBytes = pBytes[:bytesPtr]
		}
	}
	return nil
}

func (p *sConn) recvHeadBytes(
	pCtx context.Context,
	pChRead chan<- struct{},
	pInitTimeout time.Duration,
) (uint32, error) {
	defer func() { pChRead <- struct{}{} }()

	var (
		headBytes []byte
		err       error
	)

	chErr := make(chan error)
	go func() {
		headBytes, err = p.recvDataBytes(pCtx, encoding.CSizeUint32, pInitTimeout)
		if err != nil {
			chErr <- errors.Join(ErrReadHeaderBlock, err)
			return
		}
		chErr <- nil
	}()

	select {
	case <-pCtx.Done():
		return 0, pCtx.Err()
	case err := <-chErr:
		if err != nil {
			return 0, err
		}
	}

	msgSizeBytes := [encoding.CSizeUint32]byte{}
	copy(msgSizeBytes[:], headBytes)

	gotMsgSize := encoding.BytesToUint32(msgSizeBytes)
	fullMsgSize := p.fSettings.GetLimitMessageSizeBytes() + layer1.CMessageHeadSize

	switch {
	case gotMsgSize < layer1.CMessageHeadSize:
		fallthrough
	case uint64(gotMsgSize) > fullMsgSize:
		return 0, ErrInvalidMsgSize
	}

	return gotMsgSize, nil
}

func (p *sConn) recvDataBytes(pCtx context.Context, pMustLen uint32, pInitTimeout time.Duration) ([]byte, error) {
	dataRaw := make([]byte, 0, pMustLen)

	if err := p.fSocket.SetReadDeadline(time.Now().Add(pInitTimeout)); err != nil {
		return nil, errors.Join(ErrSetReadDeadline, err)
	}

	mustLen := pMustLen
	for mustLen != 0 {
		select {
		case <-pCtx.Done():
			return nil, pCtx.Err()
		default:
			buffer := make([]byte, mustLen)
			n, err := p.fSocket.Read(buffer)
			if err != nil {
				return nil, errors.Join(ErrReadFromSocket, err)
			}

			dataRaw = bytes.Join(
				[][]byte{
					dataRaw,
					buffer[:n],
				},
				[]byte{},
			)

			mustLen -= uint32(n)

			if err := p.fSocket.SetReadDeadline(time.Now().Add(p.fSettings.GetReadTimeout())); err != nil {
				return nil, errors.Join(ErrSetReadDeadline, err)
			}
		}
	}

	return dataRaw, nil
}
