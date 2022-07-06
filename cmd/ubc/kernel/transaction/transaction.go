package transaction

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/crypto/hashing"

	"github.com/number571/go-peer/cmd/ubc/kernel/settings"
)

var (
	_ ITransaction = &sTransaction{}
)

type sTransaction struct {
	fPayload   []byte
	fHash      []byte
	fSign      []byte
	fValidator asymmetric.IPubKey
}

type sTransactionJSON struct {
	FPayload   []byte `json:"payload"`
	FHash      []byte `json:"hash"`
	FSign      []byte `json:"sign"`
	FValidator []byte `json:"validator"`
}

func NewTransaction(priv asymmetric.IPrivKey, payLoad []byte) ITransaction {
	if priv == nil {
		return nil
	}

	keySize := settings.GSettings.Get(settings.CSizeAkey).(uint64)
	if priv.Size() != keySize {
		return nil
	}

	payloadSize := settings.GSettings.Get(settings.CSizePayl).(uint64)
	if uint64(len(payLoad)) > payloadSize {
		return nil
	}

	tx := &sTransaction{
		fPayload:   payLoad,
		fValidator: priv.PubKey(),
	}

	tx.fHash = tx.newHash()
	tx.fSign = priv.Sign(tx.fHash)

	return tx
}

func LoadTransaction(rawTX interface{}) ITransaction {
	switch x := rawTX.(type) {
	case []byte:
		txJSON := new(sTransactionJSON)
		err := json.Unmarshal(x, txJSON)
		if err != nil {
			return nil
		}

		tx := &sTransaction{
			fPayload:   txJSON.FPayload,
			fHash:      txJSON.FHash,
			fSign:      txJSON.FSign,
			fValidator: asymmetric.LoadRSAPubKey(txJSON.FValidator),
		}

		if !tx.IsValid() {
			return nil
		}

		return tx
	case string:
		var (
			prefix = "TX{"
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
		return LoadTransaction(pbytes)
	default:
		panic("unsupported type")
	}
}

func (tx *sTransaction) Payload() []byte {
	return tx.fPayload
}

func (tx *sTransaction) Hash() []byte {
	return tx.fHash
}

func (tx *sTransaction) Sign() []byte {
	return tx.fSign
}

func (tx *sTransaction) Validator() asymmetric.IPubKey {
	return tx.fValidator
}

func (tx *sTransaction) Bytes() []byte {
	txJSON := &sTransactionJSON{
		FPayload:   tx.Payload(),
		FHash:      tx.Hash(),
		FSign:      tx.Sign(),
		FValidator: tx.Validator().Bytes(),
	}

	txbytes, err := json.Marshal(txJSON)
	if err != nil {
		return nil
	}

	return txbytes
}

func (tx *sTransaction) String() string {
	return fmt.Sprintf("TX{%X}", tx.Bytes())
}

func (tx *sTransaction) IsValid() bool {
	payloadSize := settings.GSettings.Get(settings.CSizePayl).(uint64)
	if uint64(len(tx.fPayload)) > payloadSize {
		return false
	}

	if tx.Validator() == nil {
		return false
	}

	if !bytes.Equal(tx.Hash(), tx.newHash()) {
		return false
	}

	return tx.Validator().Verify(tx.Hash(), tx.Sign())
}

func (tx *sTransaction) newHash() []byte {
	return hashing.NewSHA256Hasher(bytes.Join(
		[][]byte{
			tx.Validator().Bytes(),
			tx.Payload(),
		},
		[]byte{},
	)).Bytes()
}
