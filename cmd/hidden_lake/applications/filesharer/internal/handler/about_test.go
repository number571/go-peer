// nolint: goerr113
package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/config"
	std_logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/logger"
)

func TestAboutPage(t *testing.T) {
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

	handler := AboutPage(httpLogger, &config.SConfig{
		FSettings: &config.SConfigSettings{
			FLanguage: "ENG",
		},
	})

	if err := aboutRequestOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := aboutRequest404(handler); err == nil {
		t.Error("request success with invalid path")
		return
	}
}

func aboutRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/about", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func aboutRequest404(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/about/undefined", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}
