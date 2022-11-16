package handler

import (
	"fmt"
	"testing"

	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/hlc"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/settings/testutils"
)

func TestHandleFriendsAPI(t *testing.T) {
	wcfg, node, srv := testAllCreate(tcPathConfig, tcPathDB, testutils.TgAddrs[7])
	defer testAllFree(node, srv)

	client := hlc.NewClient(
		hlc.NewRequester(fmt.Sprintf("http://%s", testutils.TgAddrs[7])),
	)

	aliasName := "test_name4"
	testGetFriends(t, client, wcfg.Config())
	testAddFriend(t, client, aliasName)
	testDelFriend(t, client, aliasName)
}

func testGetFriends(t *testing.T, client hlc.IClient, cfg config.IConfig) {
	friends, err := client.GetFriends()
	if err != nil {
		t.Error(err)
		return
	}

	if len(friends) != 3 {
		t.Errorf("length of friends != 3")
		return
	}

	for k, v := range friends {
		v1, ok := cfg.Friends()[k]
		if !ok {
			t.Errorf("undefined friend '%s'", k)
			return
		}
		if v.String() != v1.String() {
			t.Errorf("public keys not equals for '%s'", k)
			return
		}
	}
}

func testAddFriend(t *testing.T, client hlc.IClient, aliasName string) {
	err := client.AddFriend(
		aliasName,
		asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0]),
	)
	if err != nil {
		t.Error(err)
		return
	}

	friends, err := client.GetFriends()
	if err != nil {
		t.Error(err)
		return
	}

	if _, ok := friends[aliasName]; !ok {
		t.Errorf("undefined new public key by '%s'", aliasName)
		return
	}
}

func testDelFriend(t *testing.T, client hlc.IClient, aliasName string) {
	err := client.DelFriend(aliasName)
	if err != nil {
		t.Error(err)
		return
	}

	friends, err := client.GetFriends()
	if err != nil {
		t.Error(err)
		return
	}

	if _, ok := friends[aliasName]; ok {
		t.Errorf("deleted public key exists for '%s'", aliasName)
		return
	}
}
