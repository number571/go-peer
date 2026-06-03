package logger

import (
	"bytes"
	"testing"
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

	builder := NewLogBuilder(tcService).
		WithHash([]byte(tcHash)).
		WithProof(tcProof).
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
