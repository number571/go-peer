package chain

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/cmd/ubc/kernel/block"
	"github.com/number571/go-peer/cmd/ubc/kernel/settings"
	"github.com/number571/go-peer/cmd/ubc/kernel/transaction"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/crypto/random"
)

const (
	numBlocks = 10
	tcPrivKey = "Priv(go-peer/rsa){3082025C02010002818100C6498F4C2E238D3FA4D133E1A53A2C76D50752EAE6B9F716E0C7E1ABE47B49F0F141F1BD9411DB5A5285E9C14A20B0E703885D06513A79950C7EE6E31EEE24CECF044D62105C6CA6DD134C7E51513B16B1454BFBBB6F62A26111C3C89FD091E1A094985D50F3E9C3DF4AFE7E22C95AC62E52B9F7677E41EAA2BE0FBA0E60D731020301000102818060541017743AB53DFBF5DDFC7AE65DFF84D24007F9FAD1FCFD4A5D69C25FDAB6009E86B010A4F42956F9D36BA1756C3B6E4DEAD34CD6D985FD42112CB933FC10D7AFDDD125A907C3619BF2EAE809DDCE935E3E67AAE84A43E8E0330074F957F1024D803A8444CCDE6160BF189CCA9401ACA5DF509CD1D6B2754ABB8AC5D1D499024100FDC116B09235544664EFC0740B0A84BCCF9C1639CD323586ED6EFF39623458506E40944A1BBD2F4E29D92E78636AB17A38FCD7ED21D5A4F48A0A5B1FA5CDBA83024100C80ACDD59989B60922DD50B57C51D3681AA59C0D5786398120136F6736D7720295ACD02E1FC8627809E1071F5DBADD9FA40810B877024558EED51DAEDE82C93B02402BC5A80D535B41AB56F40885BBF5D789DE62356F491735268E448C6030B188DE6EF652DE29C4CBA9370CD0B851A5F0F17D6D182E3E9CE4F48DEF5562B32E36D3024100BDEE042099F6B66F563AEB3665230BA5FC26E15389965762D221A1D44DADA101F33A712E59DED81F40C1F71140DCFB2F677E80E1A39CF45ACBE86C966B8DA1A1024032E37D5ED53666986AEAB0FCB882CB4235B7820DBA4F8A8D5D77946C34998FA80156812477A82C7F24493CCE2066E6973E4E4375539CA34463B3B8BF7DBCCE43}"
)

func TestChain(t *testing.T) {
	const (
		chainName = "chain.db"
	)

	defer os.RemoveAll(chainName)

	priv := asymmetric.LoadRSAPrivKey(tcPrivKey)
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
			t.Error(fmt.Sprintf("failed accept block (%d)", i))
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
				t.Error(fmt.Sprintf("failed get saved transaction in chain (%d)", i))
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
