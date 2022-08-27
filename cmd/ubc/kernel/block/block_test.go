package block

import (
	"fmt"
	"testing"

	"github.com/number571/go-peer/cmd/ubc/kernel/settings"
	"github.com/number571/go-peer/cmd/ubc/kernel/transaction"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/crypto/hashing"
	"github.com/number571/go-peer/settings/testutils"
)

func TestTransaction(t *testing.T) {
	priv := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey)
	hash := hashing.NewSHA256Hasher([]byte("prev-hash")).Bytes()

	txSize := settings.GSettings.Get(settings.CSizeTrns).(uint64)
	txs := []transaction.ITransaction{}
	for i := uint64(0); i < txSize; i++ {
		tx := transaction.NewTransaction(priv,
			[]byte(fmt.Sprintf("transaction-%d", i)))
		txs = append(txs, tx)
	}

	newBlock := NewBlock(priv, hash, txs)
	if newBlock == nil {
		t.Errorf("new block is nil")
		return
	}

	if !newBlock.IsValid() {
		t.Errorf("new block is not valid")
		return
	}

	loadBlock := LoadBlock(testutils.TcLargeBody)
	if loadBlock == nil {
		t.Errorf("load block is nil")
		return
	}

	if !loadBlock.IsValid() {
		t.Errorf("load block is not valid")
		return
	}
}
