package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/pkg/settings"
	"github.com/number571/go-peer/pkg/logger"
)

func testRunService(addr string, innerAddr string) *http.Server {
	mux := http.NewServeMux()

	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{},
		FServices: map[string]string{
			tcServiceHost: innerAddr,
		},
	}

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	mux.HandleFunc(settings.CHandleIndexPath, HandleIndexAPI(logger))
	mux.HandleFunc(settings.CHandleConfigSettings, HandleConfigSettingsAPI(cfg, logger))
	mux.HandleFunc(settings.CHandleNetworkDistributePath, HandleNetworkDistributeAPI(context.Background(), cfg, logger))

	srv := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: time.Second,
	}

	go func() { _ = srv.ListenAndServe() }()
	return srv
}

func testRunInnerService(addr string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(tcServicePath, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(tcServiceHeadKey) != tcServiceHeadValue {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		reqBytes, _ := io.ReadAll(r.Body)
		fmt.Fprintf(w, "%s: %s", tcResponseMsg, string(reqBytes))
	})

	srv := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: time.Second,
	}

	go func() { _ = srv.ListenAndServe() }()
	return srv
}
