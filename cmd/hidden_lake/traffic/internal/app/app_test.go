package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/conn_keeper"

	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	testutils "github.com/number571/go-peer/test/_data"
	anon_testutils "github.com/number571/go-peer/test/_data/anonymity"
)

const (
	tcPathDB     = hlt_settings.CPathDB
	tcPathConfig = hlt_settings.CPathCFG
)

func testDeleteFiles() {
	os.RemoveAll(tcPathDB)
	os.RemoveAll(tcPathConfig)
}

func TestApp(t *testing.T) {
	testDeleteFiles()
	defer testDeleteFiles()

	cfg, err := config.NewConfig(
		tcPathConfig,
		&config.SConfig{
			FAddress: testutils.TgAddrs[23],
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	db := database.NewKeyValueDB(
		database.NewSettings(&database.SSettings{FPath: tcPathDB}),
	)

	sett := anonymity.NewSettings(&anonymity.SSettings{})
	connKeeper := conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(&conn_keeper.SSettings{}),
		anon_testutils.TestNewNetworkNode().Handle(
			sett.GetNetworkMask(), // default value
			func(_ network.INode, _ conn.IConn, _ []byte) {
				// pass response actions
			},
		),
	)

	app := NewApp(cfg, db, connKeeper)
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

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			fmt.Sprintf("http://%s", testutils.TgAddrs[23]),
			message.NewParams(
				anon_testutils.TCMessageSize,
				anon_testutils.TCWorkSize,
			),
		),
	)

	title, err := hltClient.GetIndex()
	if err != nil {
		t.Error(err)
		return
	}

	if title != hlt_settings.CTitlePattern {
		t.Error("title is incorrect")
		return
	}
}
