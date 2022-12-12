package mempool

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/cmd/ubc/kernel/transaction"
	"github.com/number571/go-peer/internal/testutils"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/encoding"
)

const (
	tcTX     = "TX{7B227061796C6F6164223A22614756736247387349486476636D786B49513D3D222C2268617368223A22467077543567634D58306D454E615A51594B7A426237397177685767465A6A65703051574F7448333474383D222C227369676E223A224436572F326C76797066766F42712B486F584A2B527A79427A444B584542546855345277335270524B32584E4B316B657A622B5A525649712B325271312F6B7354766F6A6A4548366A546A374E736F527649485A374B456C50766B483539537468487678624F6870487A695262526D7176744C4261726E5155576B6645346F6B6B36394F6D554E64764F326E7853314E497235475536756179695430684A38534147764A4C58664F304D733D222C2276616C696461746F72223A224D49474A416F4742414D5A4A6A3077754934302F704E457A346155364C48625642314C7135726E33467544483461766B65306E7738554878765A51523231705368656E425369437735774F4958515A524F6E6D564448376D347837754A4D375042453169454678737074305454483552555473577355564C2B37747659714A68456350496E394352346143556D463151382B6E443330722B66694C4A5773597555726E335A33354236714B2B44376F4F594E637841674D424141453D227D}"
	tcHashTX = "169c13e6070c5f498435a65060acc16fbf6ac215a01598dea744163ad1f7e2df"
)

func TestMempool(t *testing.T) {
	const (
		mempoolName = "mempool.db"
	)

	os.RemoveAll(mempoolName)
	defer os.RemoveAll(mempoolName)

	sett := NewSettings(&SSettings{})

	mempool := NewMempool(sett, mempoolName)
	defer mempool.Close()

	if mempool.Height() != 0 {
		t.Errorf("init mempool with height != 0")
		return
	}

	tx := transaction.LoadTransaction(
		sett.GetBlockSettings().GetTransactionSettings(),
		tcTX,
	)

	mempool.Push(tx)
	if mempool.Height() != 1 {
		t.Errorf("mempool with 1 push has height != 1")
		return
	}

	loadTX := mempool.Transaction(encoding.HexDecode(tcHashTX))
	if loadTX == nil {
		t.Errorf("load tx from mempool = nil")
		return
	}

	priv := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey1024)
	blockCountTXs := sett.GetBlockSettings().GetCountTXs()
	for i := uint64(0); i < blockCountTXs-1; i++ {
		tx := transaction.NewTransaction(
			sett.GetBlockSettings().GetTransactionSettings(),
			priv,
			[]byte(fmt.Sprintf("transaction-%d", i)),
		)
		mempool.Push(tx)
	}

	if mempool.Height() != blockCountTXs {
		t.Errorf("mempool height != blockCountTXs")
		return
	}

	txs := mempool.Pop()
	if len(txs) != int(blockCountTXs) {
		fmt.Println(len(txs), blockCountTXs)
		t.Errorf("len of pop txs != blockCountTXs")
		return
	}
}
