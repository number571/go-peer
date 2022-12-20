package hmc

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestBuilder(t *testing.T) {
	client := client.NewClient(
		client.NewSettings(&client.SSettings{
			FWorkSize:    10,
			FMessageSize: (100 << 10),
		}),
		asymmetric.LoadRSAPrivKey(testutils.TcPrivKey),
	)

	hashRecv := client.PubKey().Address().Bytes()
	builder := NewBuilder(client)

	bSize := builder.Size()
	if !bytes.Equal(bSize.FReceiver, hashRecv) {
		t.Error("builder size error (hash receiver)")
		return
	}

	bLoad := builder.Load(1)
	if bLoad.FIndex != 1 || !bytes.Equal(bLoad.FReceiver, hashRecv) {
		t.Error("builder load error (index, hash receiver)")
		return
	}

	pl := payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcBody))
	bPush := builder.Push(client.PubKey(), pl)
	if !bytes.Equal(bPush.FReceiver, hashRecv) {
		t.Error("builder push error (hash receiver)")
		return
	}

	msg := message.LoadMessage(bPush.FPackage)
	if msg == nil {
		t.Error("builder push error (message is nil [1])")
		return
	}

	pubKey, pl, err := client.Decrypt(msg)
	if err != nil {
		t.Error("builder push error (message is nil [2])")
		return
	}

	if pl.Head() != uint64(testutils.TcHead) {
		t.Error("builder push error (head is not equal)")
		return
	}

	if string(pl.Body()) != testutils.TcBody {
		t.Error("builder push error (body is not equal)")
		return
	}

	if pubKey.Address().String() != client.PubKey().Address().String() {
		t.Error("builder push error (public key is not equal)")
		return
	}
}
