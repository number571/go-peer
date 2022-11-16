package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/hlc"
	anon_testutils "github.com/number571/go-peer/modules/network/anonymity/testutils"
	"github.com/number571/go-peer/settings/testutils"
)

const (
	tcPathDB     = "database_test.db"
	tcPathConfig = "config_test.cfg"
)

func TestApp(t *testing.T) {
	defer func() {
		os.RemoveAll(tcPathDB)
		os.RemoveAll(tcPathConfig)
	}()

	cfg, err := config.NewConfig(tcPathConfig, &config.SConfig{
		FAddress: &config.SAddress{
			FTCP:  testutils.TgAddrs[14],
			FHTTP: testutils.TgAddrs[15],
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	node := anon_testutils.TestNewNode(tcPathDB)
	app := NewApp(cfg, node)
	if err := app.Run(); err != nil {
		t.Error(err)
		return
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
			return
		}
	}()

	client := hlc.NewClient(
		hlc.NewRequester(fmt.Sprintf("http://%s", testutils.TgAddrs[15])),
	)
	pubKey, err := client.PubKey()
	if err != nil {
		t.Error(err)
		return
	}

	if pubKey.String() != node.Queue().Client().PubKey().String() {
		t.Errorf("public keys are not equals")
		return
	}
}
