package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/number571/go-peer/cmd/hlm/config"
	hlm_settings "github.com/number571/go-peer/cmd/hlm/settings"
	"github.com/number571/go-peer/cmd/hls/hlc"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/logger"
)

func hlmDefaultInit() error {
	var (
		initOnly bool
		readOnly string
	)

	flag.BoolVar(&initOnly, "init", false, "run initialization only")
	flag.StringVar(&readOnly, "read-only", "", "read-only mode from one channel")
	flag.Parse()

	gConfig = config.NewConfig(hlm_settings.CPathCFG)
	if initOnly {
		os.Exit(0)
	}

	gLogger = logger.NewLogger(os.Stdout, os.Stdout, os.Stdout)
	if readOnly != "" {
		gChannelPubKey = asymmetric.LoadRSAPubKey(readOnly)
		if gChannelPubKey == nil {
			return fmt.Errorf("public key is invalid")
		}
	}

	gActions = newActions()
	gClient = hlc.NewClient(
		hlc.NewRequester(fmt.Sprintf("http://%s", gConfig.Connection())),
	)
	if _, err := gClient.PubKey(); err != nil {
		return err
	}

	gServerHTTP = initServerHTTP(gConfig, gLogger)
	return nil
}

func initServerHTTP(cfg config.IConfig, logger logger.ILogger) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/push", handlePushHTTP)

	srv := &http.Server{
		Addr:    cfg.Address(),
		Handler: mux,
	}

	ch := make(chan struct{})
	go func() {
		logger.Info(fmt.Sprintf("HTTP is listening [%s]...", cfg.Address()))
		close(ch)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Warning(err.Error())
		}
	}()
	<-ch
	return srv
}
