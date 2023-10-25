package http

import (
	"net/http"
	"testing"
)

const (
	tcService = "TST"
	tcFmtLog  = "service=TST method=GET path=/api/index conn=127.0.0.1:55555 message=hello_world"
)

func TestPanicLogger(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	logFunc := GetLogFunc()
	_ = logFunc("_")
}

func TestLogger(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/index", nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.RemoteAddr = "127.0.0.1:55555"

	logBuilder := NewLogBuilder(tcService, req).WithMessage("hello_world")
	logFunc := GetLogFunc()

	if logFunc(logBuilder) != tcFmtLog {
		t.Error("got invalid format")
		return
	}

	logGetter := logBuilder.Get()
	if logGetter.GetConn() != "127.0.0.1:55555" {
		t.Error("got conn != conn")
		return
	}

	if logGetter.GetMessage() != "hello_world" {
		t.Error("got message != message")
		return
	}

	if logGetter.GetMethod() != "GET" {
		t.Error("got method != method")
		return
	}

	if logGetter.GetPath() != "/api/index" {
		t.Error("got path != path")
		return
	}

	if logGetter.GetService() != "TST" {
		t.Error("got service != service")
		return
	}
}
