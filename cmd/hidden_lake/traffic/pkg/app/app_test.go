package app

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/pkg/client/message"

	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	testutils "github.com/number571/go-peer/test/_data"
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
	t.Parallel()

	testDeleteFiles()
	defer testDeleteFiles()

	cfg, err := config.BuildConfig(
		tcPathConfig,
		&config.SConfig{
			FSettings: &config.SConfigSettings{
				FMessageSizeBytes: testutils.TCMessageSize,
				FWorkSizeBits:     testutils.TCWorkSize,
				FQueuePeriodMS:    testutils.TCQueuePeriod,
				FMessagesCapacity: testutils.TCCapacity,
			},
			FNetworkKey: "_",
			FAddress: &config.SAddress{
				FHTTP: testutils.TgAddrs[23],
			},
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	app := NewApp(cfg, ".")
	if err := app.Run(); err != nil {
		t.Error(err)
		return
	}
	defer func() {
		if err := app.Stop(); err != nil {
			t.Error(err)
			return
		}
	}()

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			fmt.Sprintf("http://%s", testutils.TgAddrs[23]),
			&http.Client{Timeout: time.Minute},
			message.NewSettings(&message.SSettings{
				FMessageSizeBytes: testutils.TCMessageSize,
				FWorkSizeBits:     testutils.TCWorkSize,
			}),
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

	// try run after stop
	if err := app.Stop(); err != nil {
		t.Error(err)
		return
	}
	if err := app.Run(); err != nil {
		t.Error(err)
		return
	}
}
