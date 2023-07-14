package app

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	"github.com/number571/go-peer/internal/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	testutils "github.com/number571/go-peer/test/_data"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

const (
	tcPathDB     = pkg_settings.CPathDB
	tcPathConfig = pkg_settings.CPathCFG
)

func testDeleteFiles() {
	os.RemoveAll(tcPathDB)
	os.RemoveAll(tcPathConfig)
}

func TestApp(t *testing.T) {
	testDeleteFiles()
	defer testDeleteFiles()

	// Run application
	cfg, err := config.BuildConfig(tcPathConfig, &config.SConfig{
		SConfigSettings: settings.SConfigSettings{
			FSettings: settings.SConfigSettingsBlock{
				FMessageSizeBytes: testutils.TCMessageSize,
				FWorkSizeBits:     testutils.TCWorkSize,
				FKeySizeBits:      testutils.TcAKeySize,
				FQueuePeriodMS:    testutils.TCQueuePeriod,
			},
		},
		FNetwork: "_",
		FAddress: &config.SAddress{
			FTCP:  testutils.TgAddrs[14],
			FHTTP: testutils.TgAddrs[15],
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	privKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey)
	app := NewApp(cfg, privKey, ".")
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

	client := client.NewClient(
		client.NewBuilder(),
		client.NewRequester(
			fmt.Sprintf("http://%s", testutils.TgAddrs[15]),
			&http.Client{Timeout: time.Minute},
		),
	)

	// Check public key of node
	pubKey, err := client.GetPubKey()
	if err != nil {
		t.Error(err)
		return
	}
	if pubKey.ToString() != privKey.GetPubKey().ToString() {
		t.Errorf("public keys are not equals")
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
