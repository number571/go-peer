package mempool

import (
	"sync"

	"github.com/number571/go-peer/database"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/settings"

	ksettings "github.com/number571/go-peer/cmd/ubc/kernel/settings"
	"github.com/number571/go-peer/cmd/ubc/kernel/transaction"
)

var (
	_ IMempool = &sMempool{}
)

type sMempool struct {
	fMutex sync.Mutex
	fDB    database.IKeyValueDB
}

func NewMempool(path string) IMempool {
	mempool := &sMempool{
		fDB: database.NewLevelDB(path),
	}
	_, err := mempool.fDB.Get(getKeyHeight())
	if err != nil {
		res := encoding.Uint64ToBytes(0)
		err := mempool.fDB.Set(getKeyHeight(), res[:])
		if err != nil {
			panic(err)
		}
	}
	return mempool
}

func (mempool *sMempool) Height() uint64 {
	mempool.fMutex.Lock()
	defer mempool.fMutex.Unlock()

	return mempool.getHeight()
}

func (mempool *sMempool) Transaction(hash []byte) transaction.ITransaction {
	mempool.fMutex.Lock()
	defer mempool.fMutex.Unlock()

	return mempool.getTX(hash)
}

func (mempool *sMempool) Delete(hash []byte) {
	mempool.fMutex.Lock()
	defer mempool.fMutex.Unlock()

	mempool.deleteTX(hash)
}

func (mempool *sMempool) Close() error {
	mempool.fMutex.Lock()
	defer mempool.fMutex.Unlock()

	return mempool.fDB.Close()
}

func (mempool *sMempool) Clear() {
	mempool.fMutex.Lock()
	defer mempool.fMutex.Unlock()

	prefixTXs := ksettings.GSettings.Get(ksettings.CMaskPref).(string)
	iter := mempool.fDB.Iter([]byte(prefixTXs))
	defer iter.Close()

	// TODO: iter.Key without load transaction
	for iter.Next() {
		txBytes := iter.Value()

		tx := transaction.LoadTransaction(txBytes)
		if tx == nil {
			panic("mempool: tx is nil")
		}

		mempool.deleteTX(tx.Hash())
	}
}

func (mempool *sMempool) Push(tx transaction.ITransaction) {
	mempool.fMutex.Lock()
	defer mempool.fMutex.Unlock()

	var (
		hash      = tx.Hash()
		newHeight = uint64(mempool.getHeight() + 1)
	)

	// limit of height
	sizeMempool := ksettings.GSettings.Get(ksettings.CSizeMemp).(uint64)
	if newHeight > sizeMempool {
		return
	}

	// tx already exists
	if mempool.getTX(hash) != nil {
		return
	}

	res := encoding.Uint64ToBytes(newHeight)
	mempool.fDB.Set(getKeyHeight(), res[:])
	mempool.fDB.Set(getKeyTX(hash), tx.Bytes())
}

func (mempool *sMempool) Pop() []transaction.ITransaction {
	mempool.fMutex.Lock()
	defer mempool.fMutex.Unlock()

	// count of tx need be = block size
	sizeTXs := ksettings.GSettings.Get(ksettings.CSizeTrns).(uint64)
	if mempool.getHeight() < sizeTXs {
		return nil
	}

	var (
		txs   []transaction.ITransaction
		count uint64
	)

	sVal := ksettings.GSettings.Get(ksettings.CMaskPref).(string)
	iter := mempool.fDB.Iter([]byte(sVal))
	defer iter.Close()

	for count = 0; iter.Next() && count < sizeTXs; count++ {
		txBytes := iter.Value()

		tx := transaction.LoadTransaction(txBytes)
		if tx == nil {
			return nil
		}

		txs = append(txs, tx)
	}

	if count != sizeTXs {
		panic("count != settings.CSizeTrns")
	}

	for _, tx := range txs {
		mempool.deleteTX(tx.Hash())
	}

	return txs
}

func (mempool *sMempool) getHeight() uint64 {
	data, err := mempool.fDB.Get(getKeyHeight())
	if err != nil {
		panic("mempool: height undefined")
	}
	res := [settings.CSizeUint64]byte{}
	copy(res[:], data)
	return encoding.BytesToUint64(res)
}

func (mempool *sMempool) getTX(hash []byte) transaction.ITransaction {
	data, err := mempool.fDB.Get(getKeyTX(hash))
	if err != nil {
		return nil
	}
	return transaction.LoadTransaction(data)
}

func (mempool *sMempool) deleteTX(hash []byte) {
	var (
		newHeight = uint64(mempool.getHeight() - 1)
	)

	if mempool.getTX(hash) == nil {
		return
	}

	res := encoding.Uint64ToBytes(newHeight)
	mempool.fDB.Set(getKeyHeight(), res[:])
	mempool.fDB.Del(getKeyTX(hash))
}
