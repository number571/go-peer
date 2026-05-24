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
	tcConn    = "connection"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	pubKey := asymmetric.NewPrivKey().GetPubKey()
	builder := NewLogBuilder(tcService).
		WithHash([]byte(tcHash)).
		WithProof(tcProof).
		WithPubKey(pubKey).
		WithSize(tcSize).
		WithType(CLogInfoExist).
		WithConn(tcConn)

	getter := builder.Build()
	if getter.GetService() != tcService {
		t.Fatal("getter.GetService() != tcService")
	}

	if !bytes.Equal(getter.GetHash(), []byte(tcHash)) {
		t.Fatal("!bytes.Equal(getter.GetHash(), []byte(tcHash))")
	}

	if getter.GetProof() != tcProof {
		t.Fatal("getter.GetProof() != tcProof")
	}

	if !bytes.Equal(pubKey.ToBytes(), getter.GetPubKey().ToBytes()) {
		t.Fatal("!bytes.Equal(pubKey.ToBytes(), getter.GetPubKey().ToBytes())")
	}

	if getter.GetSize() != tcSize {
		t.Fatal("getter.GetSize() != tcSize")
	}

	if getter.GetType() != CLogInfoExist {
		t.Fatal("getter.GetType() != CLogInfoExist")
	}

	if getter.GetConn() != tcConn {
		t.Fatal("getter.GetConn() != tcConn")
		return
	}
}
