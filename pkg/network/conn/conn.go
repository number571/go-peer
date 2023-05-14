package conn

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
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
	fCipher   symmetric.ICipher
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
		fCipher:   symmetric.NewAESCipher([]byte(pSett.GetNetworkKey())),
	}
}

func (p *sConn) GetSettings() ISettings {
	return p.fSettings
}

func (p *sConn) GetSocket() net.Conn {
	return p.fSocket
}

func (p *sConn) FetchPayload(pPld payload.IPayload) (payload.IPayload, error) {
	if err := p.WritePayload(pPld); err != nil {
		return nil, err
	}

	chPld := make(chan payload.IPayload)
	go readPayload(p, chPld)

	select {
	case rpld := <-chPld:
		if rpld == nil {
			return nil, fmt.Errorf("failed: read payload")
		}
		return rpld, nil
	case <-time.After(p.fSettings.GetFetchTimeWait()):
		return nil, fmt.Errorf("failed: time out")
	}
}

func (p *sConn) Close() error {
	return p.GetSocket().Close()
}

func (p *sConn) WritePayload(pPld payload.IPayload) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	encMsgBytes := p.fCipher.EncryptBytes(
		message.NewMessage(
			pPld,
			[]byte(p.fSettings.GetNetworkKey()),
		).GetBytes(),
	)

	prng := random.NewStdPRNG()
	voidBytes := prng.GetBytes(prng.GetUint64() % p.fSettings.GetLimitVoidSize())

	// send headers with length of blocks
	if err := p.sendBlockSize(encMsgBytes); err != nil {
		return err
	}
	if err := p.sendBlockSize(voidBytes); err != nil {
		return err
	}

	// send blocks
	if err := p.sendBytes(encMsgBytes); err != nil {
		return err
	}
	if err := p.sendBytes(voidBytes); err != nil {
		return err
	}

	return nil
}

func (p *sConn) ReadPayload() payload.IPayload {
	chPld := make(chan payload.IPayload)
	go readPayload(p, chPld)
	return <-chPld
}

func (p *sConn) sendBytes(pBytes []byte) error {
	bytesPtr := uint64(len(pBytes))
	for {
		n, err := p.GetSocket().Write(pBytes[:bytesPtr])
		if err != nil {
			return err
		}

		bytesPtr = bytesPtr - uint64(n)
		pBytes = pBytes[:bytesPtr]

		if bytesPtr == 0 {
			break
		}
	}
	return nil
}

func (p *sConn) sendBlockSize(pBytes []byte) error {
	blockSize := encoding.Uint64ToBytes(uint64(len(pBytes)))
	n, err := p.GetSocket().Write(p.fCipher.EncryptBytes(blockSize[:]))
	if err != nil {
		return err
	}
	if n != symmetric.CAESBlockSize+encoding.CSizeUint64 {
		return fmt.Errorf("invalid size of sent package")
	}
	return nil
}

func (p *sConn) recvBlockSize() (uint64, error) {
	encBufLen := make([]byte, symmetric.CAESBlockSize+encoding.CSizeUint64)
	n, err := p.GetSocket().Read(encBufLen)
	if err != nil {
		return 0, err
	}
	if n != symmetric.CAESBlockSize+encoding.CSizeUint64 {
		return 0, fmt.Errorf("block size is invalid")
	}

	// mustLen = Size[u64] in uint64
	bufLen := p.fCipher.DecryptBytes(encBufLen)
	arrLen := [encoding.CSizeUint64]byte{}
	copy(arrLen[:], bufLen)

	return encoding.BytesToUint64(arrLen), nil
}

func readPayload(pConn *sConn, pChPld chan payload.IPayload) {
	var pld payload.IPayload
	defer func() {
		pChPld <- pld
	}()

	msgSize, err := pConn.recvBlockSize()
	if err != nil || msgSize > pConn.fSettings.GetMessageSize() {
		return
	}

	voidSize, err := pConn.recvBlockSize()
	if err != nil || voidSize > pConn.fSettings.GetLimitVoidSize() {
		return
	}

	mustLen := msgSize + voidSize
	dataRaw := make([]byte, 0, mustLen)
	for {
		buffer := make([]byte, mustLen)
		n, err := pConn.GetSocket().Read(buffer)
		if err != nil {
			return
		}

		dataRaw = bytes.Join(
			[][]byte{
				dataRaw,
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
		pConn.fCipher.DecryptBytes(dataRaw[:msgSize]),
		[]byte(pConn.fSettings.GetNetworkKey()),
	)
	if msg == nil {
		return
	}

	pld = msg.GetPayload()
}
