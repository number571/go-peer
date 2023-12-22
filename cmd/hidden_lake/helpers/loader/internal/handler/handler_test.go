package handler

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/settings"
	"github.com/number571/go-peer/pkg/logger"
	testutils "github.com/number571/go-peer/test/_data"
)

var (
	tgProducer = testutils.TgAddrs[42]
	tgConsumer = testutils.TgAddrs[43]
	tgTService = testutils.TgAddrs[44]
)

func testRunService(addr string) *http.Server {
	mux := http.NewServeMux()

	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessagesCapacity: testutils.TCCapacity,
			FWorkSizeBits:     testutils.TCWorkSize,
		},
		FProducers: []string{tgProducer},
		FConsumers: []string{tgConsumer},
	}

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	mux.HandleFunc(settings.CHandleIndexPath, HandleIndexAPI(logger))
	mux.HandleFunc(settings.CHandleTransferPath, HandleTransferAPI(cfg, logger))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}
