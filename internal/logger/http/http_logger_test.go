package logbuilder

import (
	"net/http"
	"testing"
)

const (
	tcService = "TST"
	tcFmtLog  = "service=TST method=GET path=/api/index conn=127.0.0.1:55555 message=hello_world"
)

func TestLogger(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/index", nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.RemoteAddr = "127.0.0.1:55555"

	logger := NewHTTPLogger(tcService, req)
	if logger.Get("hello_world") != tcFmtLog {
		t.Error("result fmtLog != tcFmtLog")
		return
	}
}
