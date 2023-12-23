package conn

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network/message"
)

const (
	// IV + Hash + PayloadHead + Proof
	cPayloadOverHeadSize = symmetric.CAESBlockSize + hashing.CSHA256Size + 2*encoding.CSizeUint64

	// IV + Uint64(encMsgSize) + Uint64(voidSize) + HMAC(encMsgSize || voidSize) + HMAC(msgBytes || voidBytes)
	cEncryptRecvHeadSize = symmetric.CAESBlockSize + 2*encoding.CSizeUint64 + 2*hashing.CSHA256Size
)

const (
	cWorkSizeKey = 1

	// first digits of PI
	cAuthSalt = "1415926535_8979323846_2643383279_5028841971_6939937510"

	// seconds digits of PI
	cCipherSalt = "5820974944_5923078164_0628620899_8628034825_3421170679"
)

var (
	_ IConn = &sConn{}
)

type sConn struct {
	fMutex    sync.Mutex
	fKeyMutex sync.Mutex

	fSocket   net.Conn
	fSettings ISettings

	fNetworkKey string
	fAuthKey    []byte
	fCipher     symmetric.ICipher
}

func NewConn(pSett ISettings, pAddr string) (IConn, error) {
	conn, err := net.Dial("tcp", pAddr)
	if err != nil {
		return nil, fmt.Errorf("tcp connect: %w", err)
	}
	return LoadConn(pSett, conn), nil
}

func LoadConn(pSett ISettings, pConn net.Conn) IConn {
	networkKey := pSett.GetNetworkKey()
	cipher, authKey := buildState(networkKey)

	return &sConn{
		fSettings:   pSett,
		fSocket:     pConn,
		fNetworkKey: networkKey,
		fAuthKey:    authKey,
		fCipher:     cipher,
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

func (p *sConn) WriteMessage(pCtx context.Context, pMsg message.IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	prng := random.NewStdPRNG()
	randVoidSize := prng.GetUint64() % (p.fSettings.GetLimitVoidSize() + 1)

	voidBytes := prng.GetBytes(randVoidSize)
	encMsgBytes := p.getCipher().EncryptBytes(pMsg.ToBytes())

	err := p.sendBytes(pCtx, bytes.Join(
		[][]byte{
			p.getHeadBytes(encMsgBytes, voidBytes),
			bytes.Join(
				[][]byte{
					encMsgBytes,
					voidBytes,
				},
				[]byte{},
			),
		},
		[]byte{},
	))
	if err != nil {
		return fmt.Errorf("send payload bytes: %w", err)
	}

	return nil
}

func (p *sConn) ReadMessage(pCtx context.Context, pChRead chan struct{}) (message.IMessage, error) {
	// large wait read deadline => the connection has not sent anything yet
	encMsgSize, voidSize, gotHash, err := p.recvHeadBytes(pCtx, pChRead, p.fSettings.GetWaitReadDeadline())
	if err != nil {
		return nil, fmt.Errorf("receive head bytes: %w", err)
	}

	dataBytes, err := p.recvDataBytes(pCtx, encMsgSize+voidSize)
	if err != nil {
		return nil, fmt.Errorf("receive data bytes: %w", err)
	}

	// check hash sum of received data
	newHash := p.getAuthData(bytes.Join(
		[][]byte{
			dataBytes[:encMsgSize],
			dataBytes[encMsgSize:],
		},
		[]byte{},
	))
	if !bytes.Equal(newHash, gotHash) {
		return nil, errors.New("got invalid hash")
	}

	// try unpack message from bytes
	msgBytes := p.getCipher().DecryptBytes(dataBytes[:encMsgSize])
	msg, err := message.LoadMessage(p.fSettings, msgBytes)
	if err != nil {
		return nil, fmt.Errorf("got invalid message bytes: %w", err)
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
			p.fSocket.SetWriteDeadline(time.Now().Add(p.fSettings.GetWriteDeadline()))

			n, err := p.fSocket.Write(pBytes[:bytesPtr])
			if err != nil {
				return fmt.Errorf("write tcp bytes: %w", err)
			}

			bytesPtr = bytesPtr - uint64(n)
			pBytes = pBytes[:bytesPtr]
		}
	}
	return nil
}

func (p *sConn) getHeadBytes(pEncMsgBytes, pVoidBytes []byte) []byte {
	encMsgSizeBytes := encoding.Uint64ToBytes(uint64(len(pEncMsgBytes)))
	voidSizeBytes := encoding.Uint64ToBytes(uint64(len(pVoidBytes)))

	return p.getCipher().EncryptBytes(bytes.Join(
		[][]byte{
			encMsgSizeBytes[:],
			voidSizeBytes[:],
			p.getAuthData(bytes.Join(
				[][]byte{
					encMsgSizeBytes[:],
					voidSizeBytes[:],
				},
				[]byte{},
			)),
			p.getAuthData(bytes.Join(
				[][]byte{
					pEncMsgBytes,
					pVoidBytes,
				},
				[]byte{},
			)),
		},
		[]byte{},
	))
}

func (p *sConn) recvHeadBytes(pCtx context.Context, pChRead chan<- struct{}, deadline time.Duration) (uint64, uint64, []byte, error) {
	defer func() { pChRead <- struct{}{} }()

	const (
		firstSizeIndex  = encoding.CSizeUint64
		secondSizeIndex = firstSizeIndex + encoding.CSizeUint64
		firstHashIndex  = secondSizeIndex + hashing.CSHA256Size
		secondHashIndex = firstHashIndex + hashing.CSHA256Size
	)

	p.fSocket.SetReadDeadline(time.Now().Add(deadline))
	encRecvHead := make([]byte, cEncryptRecvHeadSize)

	chErr := make(chan error)
	go func() {
		n, err := p.fSocket.Read(encRecvHead)
		if err != nil {
			chErr <- fmt.Errorf("read tcp header block: %w", err)
		}
		if n != cEncryptRecvHeadSize {
			chErr <- errors.New("invalid header block")
		}
		chErr <- nil
	}()

	select {
	case <-pCtx.Done():
		return 0, 0, nil, pCtx.Err()
	case err := <-chErr:
		if err != nil {
			return 0, 0, nil, err
		}
		break
	}

	recvHead := p.getCipher().DecryptBytes(encRecvHead)
	if recvHead == nil {
		return 0, 0, nil, errors.New("decrypt header bytes")
	}

	encMsgSizeBytes := [encoding.CSizeUint64]byte{}
	copy(encMsgSizeBytes[:], recvHead[:firstSizeIndex])

	voidSizeBytes := [encoding.CSizeUint64]byte{}
	copy(voidSizeBytes[:], recvHead[firstSizeIndex:secondSizeIndex])

	encMsgSize := encoding.BytesToUint64(encMsgSizeBytes)
	if encMsgSize > (p.fSettings.GetMessageSizeBytes() + cPayloadOverHeadSize) {
		return 0, 0, nil, errors.New("invalid header.encMsgSize")
	}

	voidSize := encoding.BytesToUint64(voidSizeBytes)
	if voidSize > p.fSettings.GetLimitVoidSize() {
		return 0, 0, nil, errors.New("invalid header.voidSize")
	}

	// check hash sum of received sizes
	gotHash := recvHead[secondSizeIndex:firstHashIndex]
	newHash := p.getAuthData(bytes.Join(
		[][]byte{
			encMsgSizeBytes[:],
			voidSizeBytes[:],
		},
		[]byte{},
	))
	if !bytes.Equal(newHash, gotHash) {
		return 0, 0, nil, errors.New("invalid header.auth")
	}

	return encMsgSize, voidSize, recvHead[firstHashIndex:secondHashIndex], nil
}

func (p *sConn) recvDataBytes(pCtx context.Context, pMustLen uint64) ([]byte, error) {
	dataRaw := make([]byte, 0, pMustLen)

	mustLen := pMustLen
	for mustLen != 0 {
		select {
		case <-pCtx.Done():
			return nil, pCtx.Err()
		default:
			p.fSocket.SetReadDeadline(time.Now().Add(p.fSettings.GetReadDeadline()))

			buffer := make([]byte, mustLen)
			n, err := p.fSocket.Read(buffer)
			if err != nil {
				return nil, err
			}

			dataRaw = bytes.Join(
				[][]byte{
					dataRaw,
					buffer[:n],
				},
				[]byte{},
			)

			mustLen -= uint64(n)
		}
	}

	return dataRaw, nil
}

func (p *sConn) getCipher() symmetric.ICipher {
	p.autoUpdateState()
	return p.fCipher
}

func (p *sConn) getAuthData(pData []byte) []byte {
	p.autoUpdateState()
	return hashing.NewHMACSHA256Hasher(p.fAuthKey, pData).ToBytes()
}

func (p *sConn) autoUpdateState() {
	p.fKeyMutex.Lock()
	defer p.fKeyMutex.Unlock()

	// networkKey can be updated from fSettings
	networkKey := p.fSettings.GetNetworkKey()
	if p.fNetworkKey == networkKey {
		return
	}

	// rewrite sConn fields
	p.fNetworkKey = networkKey
	p.fCipher, p.fAuthKey = buildState(p.fNetworkKey)
}

func buildState(pNetworkKey string) (symmetric.ICipher, []byte) {
	cipherKeyBuilder := keybuilder.NewKeyBuilder(cWorkSizeKey, []byte(cCipherSalt))
	authKeyBuilder := keybuilder.NewKeyBuilder(cWorkSizeKey, []byte(cAuthSalt))
	return symmetric.NewAESCipher(cipherKeyBuilder.Build(pNetworkKey)), authKeyBuilder.Build(pNetworkKey)
}
