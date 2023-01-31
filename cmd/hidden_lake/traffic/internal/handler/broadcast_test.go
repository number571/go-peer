package handler

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
	anon_testutils "github.com/number571/go-peer/test/_data/anonymity"
)

const (
	tcNodePathTemplate = "database_%s_node.db"
)

func TestBroadcastIndexAPI(t *testing.T) {
	addrNode := testutils.TgAddrs[24]
	os.RemoveAll(fmt.Sprintf(tcNodePathTemplate, addrNode))

	node := anon_testutils.TestNewNode(
		fmt.Sprintf(tcNodePathTemplate, addrNode),
	)
	defer func() {
		os.RemoveAll(fmt.Sprintf(tcNodePathTemplate, addrNode))
		node.Close()
	}()

	if err := node.Run(); err != nil {
		t.Error(err)
		return
	}

	// run node in server mode
	go func() {
		err := node.Network().Listen(addrNode)
		if err != nil {
			t.Error(err)
			return
		}
	}()

	time.Sleep(200 * time.Millisecond)

	addr := testutils.TgAddrs[23]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, connKeeper, db, hltClient := testAllRun(addr, addrNode)
	defer testAllFree(addr, srv, connKeeper, db)

	time.Sleep(200 * time.Millisecond)
	if len(connKeeper.Network().Connections()) != 1 {
		t.Error("length connections != 1")
		return
	}

	client := testNewClient()
	msg, err := client.Encrypt(
		client.PubKey(),
		payload.NewPayload(0, []byte(testutils.TcLargeBody)),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if err := hltClient.AddMessage(msg); err != nil {
		t.Error(err)
		return
	}

	if err := hltClient.DoBroadcast(); err != nil {
		t.Error(err)
		return
	}
}
