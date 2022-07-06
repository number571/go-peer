package hmc

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/local/client"
	"github.com/number571/go-peer/local/message"
	"github.com/number571/go-peer/local/payload"
	"github.com/number571/go-peer/settings"
	"github.com/number571/go-peer/testutils"
)

func TestBuilder(t *testing.T) {
	client := client.NewClient(
		settings.NewSettings(),
		asymmetric.LoadRSAPrivKey(testutils.TcPrivKey),
	)

	hashRecv := client.PubKey().Address().Bytes()
	builder := NewBuilder(client)

	bSize := builder.Size()
	if !bytes.Equal(bSize.Receiver, hashRecv) {
		t.Error("builder size error (hash receiver)")
	}

	bLoad := builder.Load(1)
	if bLoad.Index != 1 || !bytes.Equal(bLoad.Receiver, hashRecv) {
		t.Error("builder load error (index, hash receiver)")
	}

	pl := payload.NewPayload(testutils.TcHead, []byte(testutils.TcBody))
	bPush := builder.Push(client.PubKey(), pl)
	if !bytes.Equal(bPush.Receiver, hashRecv) {
		t.Error("builder push error (hash receiver)")
	}

	msg := message.LoadMessage(bPush.Package)
	if msg == nil {
		t.Error("builder push error (message is nil [1])")
	}

	pubKey, pl := client.Decrypt(msg)
	if pubKey == nil {
		t.Error("builder push error (message is nil [2])")
	}

	if pl.Head() != testutils.TcHead {
		t.Error("builder push error (head is not equal)")
	}

	if string(pl.Body()) != testutils.TcBody {
		t.Error("builder push error (body is not equal)")
	}

	if pubKey.Address().String() != client.PubKey().Address().String() {
		t.Error("builder push error (public key is not equal)")
	}
}
