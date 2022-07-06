package block

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/crypto/hashing"

	"github.com/number571/go-peer/cmd/ubc/kernel/settings"
	"github.com/number571/go-peer/cmd/ubc/kernel/transaction"
)

var (
	_ IBlock = &sBlock{}
)

type sBlock struct {
	fTXs       []transaction.ITransaction
	fPrevHash  []byte
	fCurrHash  []byte
	fSign      []byte
	fValidator asymmetric.IPubKey
}

type sBlockJSON struct {
	FTXs       [][]byte `json:"txs"`
	FPrevHash  []byte   `json:"prev_hash"`
	FCurrHash  []byte   `json:"curr_hash"`
	FSign      []byte   `json:"sign"`
	FValidator []byte   `json:"validator"`
}

func NewBlock(priv asymmetric.IPrivKey, prevHash []byte, txs []transaction.ITransaction) IBlock {
	block := &sBlock{
		fTXs:       txs,
		fPrevHash:  prevHash,
		fValidator: priv.PubKey(),
	}

	if len(prevHash) != hashing.GSHA256Size {
		return nil
	}

	if !block.txsAreValid() {
		return nil
	}

	block.fCurrHash = block.newHash()
	block.fSign = priv.Sign(block.fCurrHash)

	return block
}

func LoadBlock(rawBlock interface{}) IBlock {
	switch x := rawBlock.(type) {
	case []byte:
		blockJSON := new(sBlockJSON)
		err := json.Unmarshal(x, blockJSON)
		if err != nil {
			return nil
		}

		block := &sBlock{
			fPrevHash:  blockJSON.FPrevHash,
			fCurrHash:  blockJSON.FCurrHash,
			fSign:      blockJSON.FSign,
			fValidator: asymmetric.LoadRSAPubKey(blockJSON.FValidator),
		}

		for _, tx := range blockJSON.FTXs {
			block.fTXs = append(block.fTXs, transaction.LoadTransaction(tx))
		}

		if !block.IsValid() {
			return nil
		}

		return block
	case string:
		var (
			prefix = "Block{"
			suffix = "}"
		)

		if !strings.HasPrefix(x, prefix) {
			return nil
		}
		x = strings.TrimPrefix(x, prefix)

		if !strings.HasSuffix(x, suffix) {
			return nil
		}
		x = strings.TrimSuffix(x, suffix)

		pbytes, err := hex.DecodeString(x)
		if err != nil {
			return nil
		}
		return LoadBlock(pbytes)
	default:
		panic("unsupported type")
	}
}

func (block *sBlock) Transactions() []transaction.ITransaction {
	return block.fTXs
}

func (block *sBlock) PrevHash() []byte {
	return block.fPrevHash
}

func (block *sBlock) Bytes() []byte {
	blockJSON := &sBlockJSON{
		FPrevHash:  block.PrevHash(),
		FCurrHash:  block.Hash(),
		FSign:      block.Sign(),
		FValidator: block.Validator().Bytes(),
	}

	for _, tx := range block.fTXs {
		blockJSON.FTXs = append(blockJSON.FTXs, tx.Bytes())
	}

	blockBytes, err := json.Marshal(blockJSON)
	if err != nil {
		return nil
	}

	return blockBytes
}

func (block *sBlock) String() string {
	return fmt.Sprintf("Block{%X}", block.Bytes())
}

func (block *sBlock) Hash() []byte {
	return block.fCurrHash
}

func (block *sBlock) Sign() []byte {
	return block.fSign
}

func (block *sBlock) Validator() asymmetric.IPubKey {
	return block.fValidator
}

func (block *sBlock) IsValid() bool {
	if block.Validator() == nil {
		return false
	}

	if !block.txsAreValid() {
		return false
	}

	if !bytes.Equal(block.Hash(), block.newHash()) {
		return false
	}

	return block.Validator().Verify(block.Hash(), block.Sign())
}

func (block *sBlock) txsAreValid() bool {
	txSize := settings.GSettings.Get(settings.CSizeTrns).(uint64)
	if uint64(len(block.fTXs)) != txSize {
		return false
	}

	for _, tx := range block.fTXs {
		if tx == nil {
			return false
		}
		if !tx.IsValid() {
			return false
		}
	}

	sort.SliceStable(block.fTXs, func(i, j int) bool {
		return bytes.Compare(block.fTXs[i].Hash(), block.fTXs[j].Hash()) < 0
	})

	// TODO: uniq with signature
	for i := 0; i < len(block.fTXs)-1; i++ {
		if bytes.Equal(block.fTXs[i].Hash(), block.fTXs[i+1].Hash()) {
			return false
		}
	}

	return true
}

func (block *sBlock) newHash() []byte {
	hash := block.PrevHash()

	for _, tx := range block.fTXs {
		hash = hashing.NewSHA256Hasher(bytes.Join(
			[][]byte{
				hash,
				tx.Hash(),
			},
			[]byte{},
		)).Bytes()
	}

	return hash
}
