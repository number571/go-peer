// nolint: goerr113
package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/config"
	std_logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/logger"
)

func TestIndexPage(t *testing.T) {
	t.Parallel()

	logging, err := std_logger.LoadLogging([]string{})
	if err != nil {
		t.Error(err)
		return
	}

	httpLogger := std_logger.NewStdLogger(
		logging,
		func(_ logger.ILogArg) string {
			return ""
		},
	)

	handler := IndexPage(httpLogger, &config.SConfig{
		FSettings: &config.SConfigSettings{
			FLanguage: "ENG",
		},
	})

	if err := indexRequestOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := indexRequest404(handler); err == nil {
		t.Error("request success with invalid path")
		return
	}
}

func indexRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusFound {
		return errors.New("bad status code")
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func indexRequest404(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/undefined", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}
