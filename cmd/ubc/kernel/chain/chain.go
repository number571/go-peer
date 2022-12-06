package chain

import (
	"bytes"
	"fmt"
	"sort"
	"sync"

	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/closer"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/encoding"
	"github.com/number571/go-peer/modules/filesystem"
	"github.com/number571/go-peer/modules/storage/database"

	"github.com/number571/go-peer/cmd/ubc/kernel/block"
	"github.com/number571/go-peer/cmd/ubc/kernel/mempool"
	"github.com/number571/go-peer/cmd/ubc/kernel/transaction"
)

var (
	_ IChain = &sChain{}
)

type sChain struct {
	fMutex        sync.Mutex
	fSettings     ISettings
	fPrivKey      asymmetric.IPrivKey
	fBlocks       database.IKeyValueDB
	fTransactions database.IKeyValueDB
	fMempool      mempool.IMempool
}

func NewChain(sett ISettings, priv asymmetric.IPrivKey, genesis block.IBlock) (IChain, error) {
	if !genesis.IsValid() {
		return nil, fmt.Errorf("genesis block is invalid")
	}

	if filesystem.OpenFile(sett.GetRootPath()).IsExist() {
		return nil, fmt.Errorf("chain already exists")
	}

	chain := &sChain{
		fSettings: sett,
		fPrivKey:  priv,
		fBlocks: database.NewLevelDB(
			database.NewSettings(&database.SSettings{}),
			sett.GetBlocksPath(),
		),
		fTransactions: database.NewLevelDB(
			database.NewSettings(&database.SSettings{}),
			sett.GetTransactionsPath(),
		),
		fMempool: mempool.NewMempool(
			sett.GetMempoolSettings(),
			sett.GetMempoolPath(),
		),
	}

	if chain.fBlocks == nil {
		return nil, fmt.Errorf("chain.blocks == nil")
	}

	if chain.fTransactions == nil {
		return nil, fmt.Errorf("chain.transactions == nil")
	}

	if chain.fMempool == nil {
		return nil, fmt.Errorf("chain.mempool == nil")
	}

	chain.setHeight(0)
	chain.setBlock(genesis)

	return chain, nil
}

func LoadChain(sett ISettings, priv asymmetric.IPrivKey) (IChain, error) {
	if !filesystem.OpenFile(sett.GetRootPath()).IsExist() {
		return nil, fmt.Errorf("chain not exists")
	}

	chain := &sChain{
		fSettings: sett,
		fPrivKey:  priv,
		fBlocks: database.NewLevelDB(
			database.NewSettings(&database.SSettings{}),
			sett.GetBlocksPath(),
		),
		fTransactions: database.NewLevelDB(
			database.NewSettings(&database.SSettings{}),
			sett.GetTransactionsPath(),
		),
		fMempool: mempool.NewMempool(
			sett.GetMempoolSettings(),
			sett.GetMempoolPath(),
		),
	}

	if chain.fBlocks == nil {
		return nil, fmt.Errorf("chain.blocks == nil")
	}

	if chain.fTransactions == nil {
		return nil, fmt.Errorf("chain.transactions == nil")
	}

	if chain.fMempool == nil {
		return nil, fmt.Errorf("chain.mempool == nil")
	}

	return chain, nil
}

func (chain *sChain) Close() error {
	chain.fMutex.Lock()
	defer chain.fMutex.Unlock()

	return closer.CloseAll([]modules.ICloser{
		chain.fBlocks,
		chain.fTransactions,
		chain.fMempool,
	})
}

func (chain *sChain) Rollback(ptr uint64) bool {
	chain.fMutex.Lock()
	defer chain.fMutex.Unlock()

	if ptr > chain.getHeight() {
		return false
	}

	oldHeight := chain.getHeight()
	newHeight := oldHeight - ptr

	chain.setHeight(newHeight)
	for i := newHeight + 1; i <= oldHeight; i++ {
		chain.delBlock(i)
	}

	return true
}

func (chain *sChain) Mempool() mempool.IMempool {
	chain.fMutex.Lock()
	defer chain.fMutex.Unlock()

	return chain.getMempool()
}

func (chain *sChain) Accept(block block.IBlock) bool {
	chain.fMutex.Lock()
	defer chain.fMutex.Unlock()

	if !block.IsValid() {
		return false
	}

	lastBlock := chain.getBlock(chain.getHeight())
	if !bytes.Equal(lastBlock.Hash(), block.PrevHash()) {
		return false
	}

	for _, tx := range block.Transactions() {
		// this transaction already exists in BC
		// than break accept function
		if chain.getTransaction(tx.Hash()) != nil {
			return false
		}
	}

	mempool := chain.getMempool()
	for _, tx := range block.Transactions() {
		mempool.Delete(tx.Hash())
	}

	chain.setHeight(chain.getHeight() + 1)
	chain.setBlock(block)

	return true
}

func (chain *sChain) Merge(txs []transaction.ITransaction) bool {
	chain.fMutex.Lock()
	defer chain.fMutex.Unlock()

	var (
		lastBlock = chain.getBlock(chain.getHeight())
		resultTXs []transaction.ITransaction
	)

	resultTXs = append(resultTXs, lastBlock.Transactions()...)

	for _, tx := range txs {
		if !tx.IsValid() {
			return false
		}

		// this transaction already exists in BC
		// than pass and get another transaction
		if chain.getTransaction(tx.Hash()) != nil {
			continue
		}

		resultTXs = append(resultTXs, tx)
	}

	// nothing new transactions, all passed
	countTXs := chain.fSettings.GetMempoolSettings().GetBlockSettings().GetCountTXs()
	if uint64(len(resultTXs)) == countTXs {
		return false
	}

	// determinate function of gets slice of transactions
	sort.SliceStable(resultTXs, func(i, j int) bool {
		return bytes.Compare(resultTXs[i].Hash(), resultTXs[j].Hash()) < 0
	})
	appendTXs := resultTXs[:countTXs]
	deleteTXs := resultTXs[countTXs:]

	// create new block with appendTX transactions
	// and delete from old block deleteTX transactions
	chain.updateBlock(
		chain.getHeight(),
		block.NewBlock(
			chain.fSettings.GetMempoolSettings().GetBlockSettings(),
			chain.fPrivKey,
			lastBlock.PrevHash(),
			appendTXs,
		),
		deleteTXs,
	)
	return true
}

func (chain *sChain) Height() uint64 {
	chain.fMutex.Lock()
	defer chain.fMutex.Unlock()

	return chain.getHeight()
}

func (chain *sChain) Transaction(hash []byte) transaction.ITransaction {
	chain.fMutex.Lock()
	defer chain.fMutex.Unlock()

	return chain.getTransaction(hash)
}

func (chain *sChain) Block(height uint64) block.IBlock {
	chain.fMutex.Lock()
	defer chain.fMutex.Unlock()

	return chain.getBlock(height)
}

// Mempool

func (chain *sChain) getMempool() mempool.IMempool {
	return chain.fMempool
}

// Height

func (chain *sChain) getHeight() uint64 {
	data, err := chain.fBlocks.Get(getKeyHeight())
	if err != nil {
		panic("chain: height undefined")
	}
	res := [encoding.CSizeUint64]byte{}
	copy(res[:], data)
	return encoding.BytesToUint64(res)
}

func (chain *sChain) setHeight(height uint64) {
	res := encoding.Uint64ToBytes(height)
	chain.fBlocks.Set(getKeyHeight(), res[:])
}

// Transaction

func (chain *sChain) getTransaction(hash []byte) transaction.ITransaction {
	data, err := chain.fTransactions.Get(getKeyTX(hash))
	if err != nil {
		return nil
	}
	return transaction.LoadTransaction(
		chain.fSettings.GetMempoolSettings().GetBlockSettings().GetTransactionSettings(),
		data,
	)
}

func (chain *sChain) setTransaction(tx transaction.ITransaction) {
	chain.fTransactions.Set(getKeyTX(tx.Hash()), tx.Bytes())
}

func (chain *sChain) delTransaction(hash []byte) {
	chain.fTransactions.Del(getKeyTX(hash))
}

// Block

func (chain *sChain) getBlock(height uint64) block.IBlock {
	data, err := chain.fBlocks.Get(getKeyBlock(height))
	if err != nil {
		return nil
	}
	return block.LoadBlock(
		chain.fSettings.GetMempoolSettings().GetBlockSettings(),
		data,
	)
}

func (chain *sChain) setBlock(block block.IBlock) {
	chain.fBlocks.Set(getKeyBlock(chain.getHeight()), block.Bytes())

	for _, tx := range block.Transactions() {
		chain.setTransaction(tx)
	}
}

func (chain *sChain) delBlock(height uint64) {
	block := chain.getBlock(height)

	for _, tx := range block.Transactions() {
		chain.delTransaction(tx.Hash())
	}

	chain.fBlocks.Del(getKeyBlock(height))
}

func (chain *sChain) updateBlock(height uint64, block block.IBlock, delTXs []transaction.ITransaction) {
	mempool := chain.getMempool()
	chain.fBlocks.Set(getKeyBlock(height), block.Bytes())

	for _, tx := range block.Transactions() {
		chain.setTransaction(tx)
		mempool.Delete(tx.Hash())
	}

	for _, tx := range delTXs {
		chain.delTransaction(tx.Hash())
		mempool.Push(tx)
	}
}
