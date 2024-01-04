package app

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcPathConfig   = pkg_settings.CPathYML
	tcPathDatabase = pkg_settings.CPathDB
)

func testDeleteFiles() {
	os.RemoveAll(tcPathDatabase)
	os.RemoveAll(tcPathConfig)
}

func TestApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles()
	defer testDeleteFiles()

	// Run application
	cfg, err := config.BuildConfig(tcPathConfig, &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessagesCapacity: testutils.TCCapacity,
			FWorkSizeBits:     testutils.TCWorkSize,
		},
		FAddress: &config.SAddress{
			FInterface: testutils.TgAddrs[57],
			FIncoming:  testutils.TgAddrs[58],
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	app := NewApp(cfg, ".")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)

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
