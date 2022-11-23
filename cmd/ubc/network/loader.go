package network

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/cmd/ubc/kernel/block"
	"github.com/number571/go-peer/cmd/ubc/kernel/transaction"
	"github.com/number571/go-peer/modules/crypto/hashing"
	"github.com/number571/go-peer/modules/encoding"
	"github.com/number571/go-peer/modules/network/conn"
	"github.com/number571/go-peer/modules/payload"
)

var (
	_ ILoader = &sLoader{}
)

type sLoader struct {
	fMutex sync.Mutex
	fConn  conn.IConn
}

func NewLoader(conn conn.IConn) ILoader {
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

	rpld, err := loader.fConn.Request(payload.NewPayload(
		cMaskNetw,
		pld.ToBytes(),
	))
	if err != nil {
		return 0, err
	}

	if rpld.Head() != cMaskLoadHeight {
		return 0, fmt.Errorf("rpld.Head() != cMaskLoadHeight")
	}

	if len(rpld.Body()) != encoding.CSizeUint64 {
		return 0, fmt.Errorf("len(rpld.Body()) != utils.CSizeUint64")
	}

	res := [encoding.CSizeUint64]byte{}
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

	rpld, err := loader.fConn.Request(payload.NewPayload(
		cMaskNetw,
		pld.ToBytes(),
	))
	if err != nil {
		return nil, err
	}

	if rpld.Head() != cMaskLoadBlock {
		return nil, fmt.Errorf("rpld.Head() != cMaskLoadBlock")
	}

	return block.LoadBlock(rpld.Body()), nil
}

func (loader *sLoader) Transaction(hash []byte) (transaction.ITransaction, error) {
	loader.fMutex.Lock()
	defer loader.fMutex.Unlock()

	if len(hash) != hashing.CSHA256Size {
		return nil, fmt.Errorf("len(hash) != hashing.GSHA256Size")
	}

	pld := payload.NewPayload(
		cMaskLoadTransaction,
		hash,
	)

	rpld, err := loader.fConn.Request(payload.NewPayload(
		cMaskNetw,
		pld.ToBytes(),
	))
	if err != nil {
		return nil, err
	}

	if rpld.Head() != cMaskLoadTransaction {
		return nil, fmt.Errorf("rpld.Head() != cMaskLoadTransaction")
	}

	return transaction.LoadTransaction(rpld.Body()), nil
}
