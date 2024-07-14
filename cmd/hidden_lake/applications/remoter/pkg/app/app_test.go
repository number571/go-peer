package app

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/remoter/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/remoter/pkg/settings"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcTestdataPath = "./testdata/"
	tcPathConfig   = pkg_settings.CPathYML
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SAppError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func testDeleteFiles(prefixPath string) {
	os.RemoveAll(prefixPath + tcPathConfig)
}

func TestApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles("./")
	defer testDeleteFiles("./")

	// Run application
	cfg, err := config.BuildConfig(tcPathConfig, &config.SConfig{
		FSettings: &config.SConfigSettings{
			FExecTimeoutMS: 5000,
		},
		FAddress: &config.SAddress{
			FIncoming: testutils.TgAddrs[59],
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	app := NewApp(cfg)

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

func TestInitApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles(tcTestdataPath)
	defer testDeleteFiles(tcTestdataPath)

	if _, err := InitApp([]string{}, tcTestdataPath); err != nil {
		t.Error(err)
		return
	}
}
