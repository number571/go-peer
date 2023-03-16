package pprof

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/logger"

	_ "net/http/pprof"
)

func RunPprofService() {
	logger := logger.NewLogger(logger.NewSettings(&logger.SSettings{
		FInfo: os.Stdout,
		FWarn: os.Stdout,
		FErro: os.Stderr,
	}))
	go runPprofService(logger)
	time.Sleep(cWaitTime)
}

func runPprofService(logger logger.ILogger) {
	for i := 0; i < cRetriesNum; i++ {
		port := random.NewStdPRNG().GetUint64()%(4<<10) + 60000 // [60000;64096]
		logger.PushInfo(fmt.Sprintf("pprof running on %d port", port))

		err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.PushWarn("pprof service is closed")
			return
		}
	}

	logger.PushErro("pprof failed running")
}
