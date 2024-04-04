package handler

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestHandleFriendsAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 1)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 1)

	wcfg, node, _, cancel, srv := testAllCreate(pathCfg, pathDB, testutils.TgAddrs[7])
	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+testutils.TgAddrs[7],
			&http.Client{Timeout: time.Minute},
		),
	)

	aliasName := "test_name4"
	testGetFriends(t, client, wcfg.GetConfig())
	testAddFriend(t, client, aliasName)
	testDelFriend(t, client, aliasName)
}

func testGetFriends(t *testing.T, client hls_client.IClient, cfg config.IConfig) {
	friends, err := client.GetFriends(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if len(friends) != 3 {
		t.Error("length of friends != 3")
		return
	}

	for k, v := range friends {
		v1, ok := cfg.GetFriends()[k]
		if !ok {
			t.Errorf("undefined friend '%s'", k)
			return
		}
		if v.ToString() != v1.ToString() {
			t.Errorf("public keys not equals for '%s'", k)
			return
		}
	}
}

func testAddFriend(t *testing.T, client hls_client.IClient, aliasName string) {
	err := client.AddFriend(
		context.Background(),
		aliasName,
		asymmetric.LoadRSAPubKey(testutils.TgPubKeys[3]),
	)
	if err != nil {
		t.Error(err)
		return
	}

	friends, err := client.GetFriends(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if _, ok := friends[aliasName]; !ok {
		t.Errorf("undefined new public key by '%s'", aliasName)
		return
	}
}

func testDelFriend(t *testing.T, client hls_client.IClient, aliasName string) {
	err := client.DelFriend(context.Background(), aliasName)
	if err != nil {
		t.Error(err)
		return
	}

	friends, err := client.GetFriends(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if _, ok := friends[aliasName]; ok {
		t.Errorf("deleted public key exists for '%s'", aliasName)
		return
	}
}
