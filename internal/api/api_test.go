// nolint: goerr113
package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/utils"
)

type tsRequest tsResponse
type tsResponse struct {
	FMessage string `json:"message"`
}

type tsReadCloser struct {
	io.ReadCloser
}

func (p *tsReadCloser) Read(_ []byte) (n int, err error) { return 0, errors.New("some error") }
func (p *tsReadCloser) Close() error                     { return nil }

const (
	tcMessage          = "hello, world!"
	tcMethodNotAllowed = "method not allowed"
)

func TestLoadResponse(t *testing.T) {
	if _, err := loadResponse(0, &tsReadCloser{}); err == nil {
		t.Error("success load response with invalid readCloser")
		return
	}
}

func TestErrorsAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[41]
	unknownURL := "http://" + addr + "/unknown"

	client := &http.Client{
		Timeout: time.Minute / 4,
	}

	srv := testRunServer(addr)
	defer srv.Close()

	if _, err := Request(context.Background(), client, http.MethodGet, addr, nil); err == nil {
		t.Error("success request on incorrect url address")
		return
	}

	if _, err := Request(context.Background(), client, http.MethodGet, unknownURL, nil); err == nil {
		t.Error("success request on unknown url address")
		return
	}
}

func TestRequestResponseAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[40]
	testURL := "http://" + addr + "/test"

	client := &http.Client{
		Timeout: time.Minute / 4,
	}

	srv := testRunServer(addr)
	defer srv.Close()

	if _, err := Request(context.Background(), client, http.MethodGet, "\n\t\a", nil); err == nil {
		t.Error("success request on invalid url")
		return
	}

	if _, err := Request(context.Background(), client, http.MethodPatch, testURL, nil); err == nil {
		t.Error("PATCH: success request on method not allowed")
		return
	}

	respGET, err := Request(context.Background(), client, http.MethodGet, testURL, nil)
	if err != nil {
		t.Error(err)
		return
	}

	x := new(tsResponse)
	if err := encoding.DeserializeJSON(respGET, x); err != nil {
		t.Error(err)
		return
	}

	if x.FMessage != tcMessage {
		t.Error("GET: got message is invalid")
		return
	}

	// bytes
	respPOST1, err := Request(context.Background(), client, http.MethodPost, testURL, []byte(tcMessage))
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(respPOST1, bytes.Join([][]byte{[]byte("echo"), []byte(tcMessage)}, []byte{1})) {
		t.Error("POST1: got message is invalid")
		return
	}

	// string
	respPOST2, err := Request(context.Background(), client, http.MethodPost, testURL, tcMessage)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(respPOST2, bytes.Join([][]byte{[]byte("echo"), []byte(tcMessage)}, []byte{1})) {
		t.Error("POST2: got message is invalid")
		return
	}

	// struct
	respPOST3, err := Request(context.Background(), client, http.MethodPost, testURL, tsRequest{FMessage: tcMessage})
	if err != nil {
		t.Error(err)
		return
	}

	msg := fmt.Sprintf(`{"message":"%s"}`, tcMessage)
	if !bytes.Equal(respPOST3, bytes.Join([][]byte{[]byte("echo"), []byte(msg)}, []byte{1})) {
		t.Error("POST3: got message is invalid")
		return
	}
}

func testRunServer(addr string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			_ = Response(w, http.StatusMethodNotAllowed, tcMethodNotAllowed)
			return
		}

		if r.Method == http.MethodGet {
			_ = Response(w, http.StatusOK, tsResponse{FMessage: tcMessage})
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		_ = Response(w, http.StatusOK, bytes.Join([][]byte{[]byte("echo"), data}, []byte{1}))
	})

	srv := &http.Server{
		Addr:        addr,
		ReadTimeout: time.Second,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
	}

	go func() { _ = srv.ListenAndServe() }()
	time.Sleep(200 * time.Millisecond)
	return srv
}
