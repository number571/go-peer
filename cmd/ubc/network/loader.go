package network

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/cmd/ubc/kernel/block"
	"github.com/number571/go-peer/cmd/ubc/kernel/transaction"
	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/network"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/settings"
)

var (
	_ ILoader = &sLoader{}
)

type sLoader struct {
	fMutex sync.Mutex
	fConn  network.IConn
}

func NewLoader(conn network.IConn) ILoader {
	return &sLoader{
		fConn: conn,
	}
}

func (loader *sLoader) Height() (uint64, error) {
	loader.fMutex.Lock()
	defer loader.fMutex.Unlock()

	pld := payload.NewPayload(
		cMaskLoadHeight,
		[]byte{},
	)

	msg := loader.fConn.Request(network.NewMessage(
		payload.NewPayload(
			cMaskNetw,
			pld.Bytes(),
		),
	))

	rpld := msg.Payload()
	if rpld.Head() != cMaskLoadHeight {
		return 0, fmt.Errorf("rpld.Head() != cMaskLoadHeight")
	}

	if len(rpld.Body()) != settings.CSizeUint64 {
		return 0, fmt.Errorf("len(rpld.Body()) != settings.CSizeUint64")
	}

	res := [settings.CSizeUint64]byte{}
	copy(res[:], rpld.Body())
	return encoding.BytesToUint64(res), nil
}

func (loader *sLoader) Block(i uint64) (block.IBlock, error) {
	loader.fMutex.Lock()
	defer loader.fMutex.Unlock()

	res := encoding.Uint64ToBytes(i)
	pld := payload.NewPayload(
		cMaskLoadBlock,
		res[:],
	)

	msg := loader.fConn.Request(network.NewMessage(
		payload.NewPayload(
			cMaskNetw,
			pld.Bytes(),
		),
	))

	rpld := msg.Payload()
	if rpld.Head() != cMaskLoadBlock {
		return nil, fmt.Errorf("rpld.Head() != cMaskLoadBlock")
	}

	return block.LoadBlock(rpld.Body()), nil
}

func (loader *sLoader) Transaction(hash []byte) (transaction.ITransaction, error) {
	loader.fMutex.Lock()
	defer loader.fMutex.Unlock()

	if len(hash) != hashing.GSHA256Size {
		return nil, fmt.Errorf("len(hash) != hashing.GSHA256Size")
	}

	pld := payload.NewPayload(
		cMaskLoadTransaction,
		hash,
	)

	msg := loader.fConn.Request(network.NewMessage(
		payload.NewPayload(
			cMaskNetw,
			pld.Bytes(),
		),
	))

	rpld := msg.Payload()
	if rpld.Head() != cMaskLoadTransaction {
		return nil, fmt.Errorf("rpld.Head() != cMaskLoadTransaction")
	}

	return transaction.LoadTransaction(rpld.Body()), nil
}
