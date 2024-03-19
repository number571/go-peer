package conn

import (
	"bytes"
	"context"
	"net"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/utils"
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
	dialer := &net.Dialer{Timeout: pSett.GetDialTimeout()}
	conn, err := dialer.Dial("tcp", pAddr)
	if err != nil {
		return nil, utils.MergeErrors(ErrCreateConnection, err)
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

func (p *sConn) Close() error {
	return p.fSocket.Close()
}

func (p *sConn) WriteMessage(pCtx context.Context, pMsg net_message.IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	msgBytes := pMsg.ToBytes()
	msgSizeBytes := encoding.Uint64ToBytes(uint64(len(msgBytes)))

	err := p.sendBytes(pCtx, bytes.Join(
		[][]byte{
			msgSizeBytes[:],
			msgBytes,
		},
		[]byte{},
	))
	if err != nil {
		return utils.MergeErrors(ErrSendPayloadBytes, err)
	}

	return nil
}

func (p *sConn) ReadMessage(pCtx context.Context, pChRead chan<- struct{}) (net_message.IMessage, error) {
	// large wait read deadline => the connection has not sent anything yet
	msgSize, err := p.recvHeadBytes(pCtx, pChRead, p.fSettings.GetWaitReadTimeout())
	if err != nil {
		return nil, utils.MergeErrors(ErrReadHeaderBytes, err)
	}

	dataBytes, err := p.recvDataBytes(pCtx, msgSize, p.fSettings.GetReadTimeout())
	if err != nil {
		return nil, utils.MergeErrors(ErrReadBodyBytes, err)
	}

	// try unpack message from bytes
	msg, err := net_message.LoadMessage(p.fSettings, dataBytes)
	if err != nil {
		return nil, utils.MergeErrors(ErrInvalidMessageBytes, err)
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
			_ = p.fSocket.SetWriteDeadline(time.Now().Add(p.fSettings.GetWriteTimeout()))

			n, err := p.fSocket.Write(pBytes[:bytesPtr])
			if err != nil {
				return utils.MergeErrors(ErrWriteToSocket, err)
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
) (uint64, error) {
	defer func() { pChRead <- struct{}{} }()

	const (
		sizeIndex     = encoding.CSizeUint64
		sizeHashIndex = sizeIndex + hashing.CSHA256Size
		dataHashIndex = sizeHashIndex + hashing.CSHA256Size
	)

	sizeHead := make([]byte, encoding.CSizeUint64)
	chErr := make(chan error)

	go func() {
		var err error
		sizeHead, err = p.recvDataBytes(pCtx, encoding.CSizeUint64, pInitTimeout)
		if err != nil {
			chErr <- utils.MergeErrors(ErrReadHeaderBlock, err)
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
		break
	}

	msgSizeBytes := [encoding.CSizeUint64]byte{}
	copy(msgSizeBytes[:], sizeHead[:sizeIndex])

	gotMsgSize := encoding.BytesToUint64(msgSizeBytes)
	fullMsgSize := p.fSettings.GetLimitMessageSizeBytes() + net_message.CMessageHeadSize

	switch {
	case gotMsgSize < net_message.CMessageHeadSize:
		fallthrough
	case gotMsgSize > fullMsgSize+p.fSettings.GetLimitVoidSizeBytes():
		return 0, ErrInvalidMsgSize
	}

	return gotMsgSize, nil
}

func (p *sConn) recvDataBytes(pCtx context.Context, pMustLen uint64, pInitTimeout time.Duration) ([]byte, error) {
	dataRaw := make([]byte, 0, pMustLen)

	_ = p.fSocket.SetReadDeadline(time.Now().Add(pInitTimeout))
	mustLen := pMustLen
	for mustLen != 0 {
		select {
		case <-pCtx.Done():
			return nil, pCtx.Err()
		default:
			buffer := make([]byte, mustLen)
			n, err := p.fSocket.Read(buffer)
			if err != nil {
				return nil, utils.MergeErrors(ErrReadFromSocket, err)
			}

			dataRaw = bytes.Join(
				[][]byte{
					dataRaw,
					buffer[:n],
				},
				[]byte{},
			)

			mustLen -= uint64(n)
			_ = p.fSocket.SetReadDeadline(time.Now().Add(p.fSettings.GetReadTimeout()))
		}
	}

	return dataRaw, nil
}
