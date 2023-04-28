package block

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"

	"github.com/number571/go-peer/cmd/union_blockchain/kernel/transaction"
)

var (
	_ IBlock = &sBlock{}
)

type sBlock struct {
	fSettings  ISettings
	fTXs       []transaction.ITransaction
	fValidator asymmetric.IPubKey

	FTXs       [][]byte `json:"txs"`
	FPrevHash  []byte   `json:"prev_hash"`
	FCurrHash  []byte   `json:"curr_hash"`
	FSign      []byte   `json:"sign"`
	FValidator []byte   `json:"validator"`
}

func NewBlock(sett ISettings, priv asymmetric.IPrivKey, prevHash []byte, txs []transaction.ITransaction) IBlock {
	block := &sBlock{
		fSettings:  sett,
		fTXs:       txs,
		fValidator: priv.GetPubKey(),
		FPrevHash:  prevHash,
		FValidator: priv.GetPubKey().ToBytes(),
	}

	if len(prevHash) != hashing.CSHA256Size {
		return nil
	}

	if !block.txsAreValid() {
		return nil
	}

	for _, tx := range txs {
		block.FTXs = append(
			block.FTXs,
			tx.Bytes(),
		)
	}

	block.FCurrHash = block.newHash()
	block.FSign = priv.SignBytes(block.FCurrHash)

	return block
}

func LoadBlock(sett ISettings, rawBlock interface{}) IBlock {
	switch x := rawBlock.(type) {
	case []byte:
		block := new(sBlock)
		err := json.Unmarshal(x, block)
		if err != nil {
			return nil
		}

		block.fSettings = sett
		block.fValidator = asymmetric.LoadRSAPubKey(block.FValidator)

		for _, tx := range block.FTXs {
			block.fTXs = append(
				block.fTXs,
				transaction.LoadTransaction(
					block.fSettings.GetTransactionSettings(),
					tx,
				),
			)
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
		return LoadBlock(sett, pbytes)
	default:
		panic("unsupported type")
	}
}

func (block *sBlock) Settings() ISettings {
	return block.fSettings
}

func (block *sBlock) Transactions() []transaction.ITransaction {
	return block.fTXs
}

func (block *sBlock) PrevHash() []byte {
	return block.FPrevHash
}

func (block *sBlock) Bytes() []byte {
	blockBytes, err := json.Marshal(block)
	if err != nil {
		return nil
	}

	return blockBytes
}

func (block *sBlock) String() string {
	return fmt.Sprintf("Block{%X}", block.Bytes())
}

func (block *sBlock) Hash() []byte {
	return block.FCurrHash
}

func (block *sBlock) Sign() []byte {
	return block.FSign
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

	return block.Validator().VerifyBytes(block.Hash(), block.Sign())
}

func (block *sBlock) txsAreValid() bool {
	if uint64(len(block.fTXs)) != block.fSettings.GetCountTXs() {
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
		)).ToBytes()
	}

	return hash
}
