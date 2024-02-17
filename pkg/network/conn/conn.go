package conn

import (
	"bytes"
	"context"
	"net"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/utils"
)

const (
	cSaltSize = 16

	// IV + Proof + Salt + Hash + PayloadHead
	cEncryptMessageHeadSize = symmetric.CAESBlockSize + encoding.CSizeUint64 + message.CSaltSize + hashing.CSHA256Size + encoding.CSizeUint64

	// Salt(cipher) + Salt(auth) + IV + Uint64(encMsgSize) + Uint64(voidSize) + HMAC(encMsgSize || voidSize) + HMAC(msgBytes || voidBytes)
	cEncryptRecvHeadSize = 2*cSaltSize + symmetric.CAESBlockSize + 2*encoding.CSizeUint64 + 2*hashing.CSHA256Size
)

var (
	_ IConn = &sConn{}
)

type sConn struct {
	fMutex    sync.Mutex
	fSocket   net.Conn
	fSettings ISettings
}

type sState struct {
	fAuthKey []byte
	fCipher  symmetric.ICipher
}

func NewConn(pSett ISettings, pAddr string) (IConn, error) {
	conn, err := net.Dial("tcp", pAddr)
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

func (p *sConn) WriteMessage(pCtx context.Context, pMsg message.IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	prng := random.NewStdPRNG()

	cipherSalt, authSalt := prng.GetBytes(cSaltSize), prng.GetBytes(cSaltSize)
	state := p.buildState(cipherSalt, authSalt)

	randVoidSize := prng.GetUint64() % (p.fSettings.GetLimitVoidSize() + 1)
	voidBytes := prng.GetBytes(randVoidSize)

	encMsgBytes := state.fCipher.EncryptBytes(pMsg.ToBytes())
	err := p.sendBytes(pCtx, bytes.Join(
		[][]byte{
			p.getHeadBytes(
				state,
				cipherSalt,
				authSalt,
				encMsgBytes,
				voidBytes,
			),
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
		return utils.MergeErrors(ErrSendPayloadBytes, err)
	}

	return nil
}

func (p *sConn) ReadMessage(pCtx context.Context, pChRead chan<- struct{}) (message.IMessage, error) {
	// large wait read deadline => the connection has not sent anything yet
	encMsgSize, voidSize, state, gotHash, err := p.recvHeadBytes(pCtx, pChRead, p.fSettings.GetWaitReadDeadline())
	if err != nil {
		return nil, utils.MergeErrors(ErrReadHeaderBytes, err)
	}

	dataBytes, err := p.recvDataBytes(pCtx, encMsgSize+voidSize, p.fSettings.GetReadDeadline())
	if err != nil {
		return nil, utils.MergeErrors(ErrReadBodyBytes, err)
	}

	// check hash sum of received data
	newHash := hashing.NewHMACSHA256Hasher(
		state.fAuthKey,
		bytes.Join(
			[][]byte{
				dataBytes[:encMsgSize],
				dataBytes[encMsgSize:],
			},
			[]byte{},
		),
	).ToBytes()
	if !bytes.Equal(newHash, gotHash) {
		return nil, ErrInvalidBodyAuthHash
	}

	// try unpack message from bytes
	msgBytes := state.fCipher.DecryptBytes(dataBytes[:encMsgSize])
	msg, err := message.LoadMessage(p.fSettings, msgBytes)
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
			_ = p.fSocket.SetWriteDeadline(time.Now().Add(p.fSettings.GetWriteDeadline()))

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

func (p *sConn) getHeadBytes(
	pState *sState,
	pCipherSalt []byte,
	pAuthSalt []byte,
	pEncMsgBytes []byte,
	pVoidBytes []byte,
) []byte {
	encMsgSizeBytes := encoding.Uint64ToBytes(uint64(len(pEncMsgBytes)))
	voidSizeBytes := encoding.Uint64ToBytes(uint64(len(pVoidBytes)))

	encHeadPart := pState.fCipher.EncryptBytes(bytes.Join(
		[][]byte{
			encMsgSizeBytes[:],
			voidSizeBytes[:],
			hashing.NewHMACSHA256Hasher(
				pState.fAuthKey,
				bytes.Join(
					[][]byte{
						encMsgSizeBytes[:],
						voidSizeBytes[:],
					},
					[]byte{},
				),
			).ToBytes(),
			hashing.NewHMACSHA256Hasher(
				pState.fAuthKey,
				bytes.Join(
					[][]byte{
						pEncMsgBytes,
						pVoidBytes,
					},
					[]byte{},
				),
			).ToBytes(),
		},
		[]byte{},
	))

	return bytes.Join(
		[][]byte{
			pCipherSalt,
			pAuthSalt,
			encHeadPart,
		},
		[]byte{},
	)
}

func (p *sConn) recvHeadBytes(pCtx context.Context, pChRead chan<- struct{}, pInitDeadline time.Duration) (uint64, uint64, *sState, []byte, error) {
	defer func() { pChRead <- struct{}{} }()

	const (
		firstSizeIndex  = encoding.CSizeUint64
		secondSizeIndex = firstSizeIndex + encoding.CSizeUint64
		firstHashIndex  = secondSizeIndex + hashing.CSHA256Size
		secondHashIndex = firstHashIndex + hashing.CSHA256Size
	)

	encRecvHead := make([]byte, cEncryptRecvHeadSize)
	chErr := make(chan error)

	go func() {
		var err error
		encRecvHead, err = p.recvDataBytes(pCtx, cEncryptRecvHeadSize, pInitDeadline)
		if err != nil {
			chErr <- utils.MergeErrors(ErrReadHeaderBlock, err)
			return
		}
		chErr <- nil
	}()

	select {
	case <-pCtx.Done():
		return 0, 0, nil, nil, pCtx.Err()
	case err := <-chErr:
		if err != nil {
			return 0, 0, nil, nil, err
		}
		break
	}

	state := p.buildState(encRecvHead[:cSaltSize], encRecvHead[cSaltSize:2*cSaltSize])
	recvHead := state.fCipher.DecryptBytes(encRecvHead[2*cSaltSize:])

	encMsgSizeBytes := [encoding.CSizeUint64]byte{}
	copy(encMsgSizeBytes[:], recvHead[:firstSizeIndex])

	voidSizeBytes := [encoding.CSizeUint64]byte{}
	copy(voidSizeBytes[:], recvHead[firstSizeIndex:secondSizeIndex])

	encMsgSize := encoding.BytesToUint64(encMsgSizeBytes)
	if encMsgSize > (p.fSettings.GetMessageSizeBytes() + cEncryptMessageHeadSize) {
		return 0, 0, nil, nil, ErrInvalidHeaderMsgSize
	}

	voidSize := encoding.BytesToUint64(voidSizeBytes)
	if voidSize > p.fSettings.GetLimitVoidSize() {
		return 0, 0, nil, nil, ErrInvalidHeaderVoidSize
	}

	// check hash sum of received sizes
	gotHash := recvHead[secondSizeIndex:firstHashIndex]
	newHash := hashing.NewHMACSHA256Hasher(
		state.fAuthKey,
		bytes.Join(
			[][]byte{
				encMsgSizeBytes[:],
				voidSizeBytes[:],
			},
			[]byte{},
		),
	).ToBytes()
	if !bytes.Equal(newHash, gotHash) {
		return 0, 0, nil, nil, ErrInvalidHeaderAuthHash
	}

	return encMsgSize, voidSize, state, recvHead[firstHashIndex:secondHashIndex], nil
}

func (p *sConn) recvDataBytes(pCtx context.Context, pMustLen uint64, pInitDeadline time.Duration) ([]byte, error) {
	dataRaw := make([]byte, 0, pMustLen)

	_ = p.fSocket.SetReadDeadline(time.Now().Add(pInitDeadline))
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
			_ = p.fSocket.SetReadDeadline(time.Now().Add(p.fSettings.GetReadDeadline()))
		}
	}

	return dataRaw, nil
}

func (p *sConn) buildState(pCipherSalt, pAuthSalt []byte) *sState {
	networkKey := p.fSettings.GetNetworkKey()
	cipherKeyBuilder := keybuilder.NewKeyBuilder(1, pCipherSalt)
	authKeyBuilder := keybuilder.NewKeyBuilder(1, pAuthSalt)
	return &sState{
		fAuthKey: authKeyBuilder.Build(networkKey),
		fCipher:  symmetric.NewAESCipher(cipherKeyBuilder.Build(networkKey)),
	}
}
