package chain

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/cmd/union_blockchain/kernel/block"
	"github.com/number571/go-peer/cmd/union_blockchain/kernel/transaction"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestChain(t *testing.T) {
	const (
		numBlocks = 10
		chainName = "chain.db"
	)

	os.RemoveAll(chainName)
	defer os.RemoveAll(chainName)

	sett := NewSettings(&SSettings{FRootPath: chainName})

	priv := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey1024)
	hash := hashing.NewSHA256Hasher([]byte("prev-hash")).ToBytes()

	// generate genesis block
	chain, err := NewChain(
		sett,
		priv,
		block.NewBlock(
			sett.GetMempoolSettings().GetBlockSettings(),
			priv,
			hash,
			testGetNewTransactions(
				sett.GetMempoolSettings().GetBlockSettings(),
				priv,
			),
		),
	)
	if err != nil {
		t.Error(err)
		return
	}
	defer chain.Close()

	// generate blocks
	for i := 0; i < numBlocks; i++ {
		newBlock := block.NewBlock(
			sett.GetMempoolSettings().GetBlockSettings(),
			priv,
			testGetLastBlockHash(chain),
			testGetNewTransactions(
				sett.GetMempoolSettings().GetBlockSettings(),
				priv,
			),
		)

		// accept generated block to chain
		ok := chain.Accept(newBlock)
		if !ok {
			t.Errorf("failed accept block (%d)", i)
			return
		}

		if !bytes.Equal(
			chain.Block(chain.Height()).Hash(),
			newBlock.Hash(),
		) {
			t.Error("newBlock.Hash != lastBlock.Hash")
			return
		}

		for _, tx := range newBlock.Transactions() {
			if chain.Transaction(tx.Hash()) == nil {
				t.Errorf("failed get saved transaction in chain (%d)", i)
				return
			}
		}
	}

	if chain.Height() != numBlocks {
		t.Error("chain.Height() != numBlocks")
		return
	}

	// merge another block with last block
	lastBlock := chain.Block(chain.Height())
	anotherBlock := block.NewBlock(
		sett.GetMempoolSettings().GetBlockSettings(),
		priv,
		testGetLastBlockHash(chain),
		testGetNewTransactions(
			sett.GetMempoolSettings().GetBlockSettings(),
			priv,
		),
	)

	ok := chain.Merge(anotherBlock.Transactions())
	if !ok {
		t.Error("failed merge transactions")
		return
	}

	txsEqual := true
	for i := range lastBlock.Transactions() {
		if !bytes.Equal(
			lastBlock.Transactions()[i].Hash(),
			anotherBlock.Transactions()[i].Hash(),
		) {
			txsEqual = false
			break
		}
	}

	if txsEqual {
		t.Error("failed merge transactions as result")
		return
	}

	// rollback chain by one block
	preLastBlock := chain.Block(chain.Height() - 1)

	ok = chain.Rollback(1)
	if !ok {
		t.Error("failed rollback chain")
		return
	}

	if !bytes.Equal(
		preLastBlock.Hash(),
		testGetLastBlockHash(chain),
	) {
		t.Error("failed rollback chain as result")
		return
	}
}

func testGetLastBlockHash(chain IChain) []byte {
	return chain.Block(chain.Height()).Hash()
}

func testGetNewTransactions(sett block.ISettings, priv asymmetric.IPrivKey) []transaction.ITransaction {
	txs := []transaction.ITransaction{}
	for i := uint64(0); i < sett.GetCountTXs(); i++ {
		txs = append(
			txs,
			transaction.NewTransaction(
				sett.GetTransactionSettings(),
				priv,
				[]byte(fmt.Sprintf("transaction-%d-%X", i, random.NewStdPRNG().GetBytes(20))),
			),
		)
	}
	return txs
}
