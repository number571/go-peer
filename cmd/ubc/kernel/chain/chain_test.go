package chain

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/cmd/ubc/kernel/block"
	"github.com/number571/go-peer/cmd/ubc/kernel/settings"
	"github.com/number571/go-peer/cmd/ubc/kernel/transaction"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/crypto/hashing"
	"github.com/number571/go-peer/modules/crypto/random"
	"github.com/number571/go-peer/settings/testutils"
)

const (
	numBlocks = 10
)

func TestChain(t *testing.T) {
	const (
		chainName = "chain.db"
	)

	os.RemoveAll(chainName)
	defer os.RemoveAll(chainName)

	priv := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey)
	hash := hashing.NewSHA256Hasher([]byte("prev-hash")).Bytes()

	// generate genesis block
	chain, err := NewChain(
		priv,
		chainName,
		block.NewBlock(priv, hash, testGetNewTransactions(priv)),
	)
	if err != nil {
		t.Error(err)
		return
	}
	defer chain.Close()

	// generate blocks
	for i := 0; i < numBlocks; i++ {
		newBlock := block.NewBlock(
			priv,
			testGetLastBlockHash(chain),
			testGetNewTransactions(priv),
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
		priv,
		testGetLastBlockHash(chain),
		testGetNewTransactions(priv),
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

func testGetNewTransactions(priv asymmetric.IPrivKey) []transaction.ITransaction {
	txSize := settings.GSettings.Get(settings.CSizeTrns).(uint64)
	txs := []transaction.ITransaction{}
	for i := uint64(0); i < txSize; i++ {
		txs = append(
			txs,
			transaction.NewTransaction(
				priv,
				[]byte(fmt.Sprintf("transaction-%d-%X", i, random.NewStdPRNG().Bytes(20))),
			),
		)
	}
	return txs
}
