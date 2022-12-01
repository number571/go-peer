package transaction

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/crypto/hashing"
)

var (
	_ ITransaction = &sTransaction{}
)

type sTransaction struct {
	fSettings  ISettings
	fValidator asymmetric.IPubKey

	FPayload   []byte `json:"payload"`
	FHash      []byte `json:"hash"`
	FSign      []byte `json:"sign"`
	FValidator []byte `json:"validator"`
}

func NewTransaction(sett ISettings, priv asymmetric.IPrivKey, payLoad []byte) ITransaction {
	tx := &sTransaction{
		fSettings:  sett,
		fValidator: priv.PubKey(),
		FPayload:   payLoad,
		FValidator: priv.PubKey().Bytes(),
	}

	tx.FHash = tx.newHash()
	tx.FSign = priv.Sign(tx.FHash)

	if !tx.IsValid() {
		return nil
	}

	return tx
}

func LoadTransaction(sett ISettings, rawTX interface{}) ITransaction {
	switch x := rawTX.(type) {
	case []byte:
		tx := new(sTransaction)
		err := json.Unmarshal(x, tx)
		if err != nil {
			return nil
		}

		tx.fSettings = sett
		tx.fValidator = asymmetric.LoadRSAPubKey(tx.FValidator)

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
		return LoadTransaction(sett, pbytes)
	default:
		panic("unsupported type")
	}
}

func (tx *sTransaction) Settings() ISettings {
	return tx.fSettings
}

func (tx *sTransaction) Payload() []byte {
	return tx.FPayload
}

func (tx *sTransaction) Hash() []byte {
	return tx.FHash
}

func (tx *sTransaction) Sign() []byte {
	return tx.FSign
}

func (tx *sTransaction) Validator() asymmetric.IPubKey {
	return tx.fValidator
}

func (tx *sTransaction) Bytes() []byte {
	txbytes, err := json.Marshal(tx)
	if err != nil {
		return nil
	}

	return txbytes
}

func (tx *sTransaction) String() string {
	return fmt.Sprintf("TX{%X}", tx.Bytes())
}

func (tx *sTransaction) IsValid() bool {
	if uint64(len(tx.FPayload)) > tx.fSettings.GetMaxSize() {
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
