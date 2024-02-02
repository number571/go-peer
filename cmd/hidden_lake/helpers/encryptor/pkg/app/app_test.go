package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/client"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcPathConfig = pkg_settings.CPathYML
)

func testDeleteFiles() {
	os.RemoveAll(tcPathConfig)
}

func TestApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles()
	defer testDeleteFiles()

	// Run application
	cfg, err := config.BuildConfig(tcPathConfig, &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FWorkSizeBits:     testutils.TCWorkSize,
			FKeySizeBits:      testutils.TcKeySize,
			FNetworkKey:       "_",
		},
		FAddress: &config.SAddress{
			FHTTP: testutils.TgAddrs[55],
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	privKey := asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024)
	app := NewApp(cfg, privKey, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)
	client := client.NewClient(
		client.NewRequester(
			fmt.Sprintf("http://%s", testutils.TgAddrs[55]),
			&http.Client{Timeout: time.Minute},
			testNetworkMessageSettings(),
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

	// try twice running
	go func() {
		if err := app.Run(ctx); err == nil {
			t.Error("success double run")
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	// try twice running
	go func() {
		if err := app.Run(ctx1); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
			return
		}
	}()
	time.Sleep(100 * time.Millisecond)
}

func testNetworkMessageSettings() net_message.ISettings {
	return net_message.NewSettings(&net_message.SSettings{
		FNetworkKey:   testutils.TCNetworkKey,
		FWorkSizeBits: testutils.TCWorkSize,
	})
}
