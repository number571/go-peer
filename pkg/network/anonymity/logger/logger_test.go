package logger

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"

	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcService = "ServiceName"
	tcRecv    = true
	tcHash    = "hash-example"
	tcProof   = 3
	tcSize    = 8192
)

func TestLogger(t *testing.T) {
	t.Parallel()

	pubKey := asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0])
	builder := NewLogBuilder(tcService).
		WithConn(nil).
		WithRecv(tcRecv).
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

	if getter.GetRecv() != tcRecv {
		t.Error("getter.GetRecv() != tcRecv")
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

	if getter.GetPubKey().GetHasher().ToString() != pubKey.GetHasher().ToString() {
		t.Error("getter.GetPubKey().GetHasher().ToString() != pubKey.GetHasher().ToString()")
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
