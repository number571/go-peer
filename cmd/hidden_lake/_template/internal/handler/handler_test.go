package handler

import (
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/_template/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/_template/pkg/settings"
	"github.com/number571/go-peer/pkg/logger"
)

func testRunService(addr string) *http.Server {
	mux := http.NewServeMux()

	// TODO: need implementation
	_ = &config.SConfig{}

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	mux.HandleFunc(settings.CHandleIndexPath, HandleIndexAPI(logger))

	srv := &http.Server{
		Addr:        addr,
		ReadTimeout: time.Second,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}
