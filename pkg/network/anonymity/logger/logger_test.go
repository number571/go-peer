package logger

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

const (
	tcService = "ServiceName"
	tcHash    = "hash-example"
	tcProof   = 3
	tcSize    = 8192
)

func TestLogger(t *testing.T) {
	t.Parallel()

	pubKey := asymmetric.NewPrivKey().GetPubKey().GetDSAPubKey()
	builder := NewLogBuilder(tcService).
		WithConn(nil).
		WithHash([]byte(tcHash)).
		WithProof(tcProof).
		WithPubKey(pubKey).
		WithSize(tcSize).
		WithType(CLogInfoExist)

	getter := builder.Get()
	if getter.GetService() != tcService {
		t.Error("getter.GetService() != tcService")
		return
	}

	if getter.GetConn() != nil {
		t.Error("getter.GetConn() != nil")
		return
	}

	if !bytes.Equal(getter.GetHash(), []byte(tcHash)) {
		t.Error("!bytes.Equal(getter.GetHash(), []byte(tcHash))")
		return
	}

	if getter.GetProof() != tcProof {
		t.Error("getter.GetProof() != tcProof")
		return
	}

	if !bytes.Equal(pubKey.ToBytes(), getter.GetPubKey().ToBytes()) {
		t.Error("!bytes.Equal(pubKey.ToBytes(), getter.GetPubKey().ToBytes())")
		return
	}

	if getter.GetSize() != tcSize {
		t.Error("getter.GetSize() != tcSize")
		return
	}

	if getter.GetType() != CLogInfoExist {
		t.Error("getter.GetType() != CLogInfoExist")
		return
	}
}
