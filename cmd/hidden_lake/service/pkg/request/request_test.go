package request

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/number571/go-peer/pkg/payload/joiner"
)

const (
	tcHost   = "test_host"
	tcPath   = "test_path"
	tcMethod = "test_method"
)

var (
	tgHead = map[string]string{
		"test_header1": "test_value1",
		"test_header2": "test_value2",
		"test_header3": "test_value3",
	}
	tgBody     = []byte("test_data")
	tgBRequest = `{"method":"test_method","host":"test_host","path":"test_path","head":{"test_header1":"test_value1","test_header2":"test_value2","test_header3":"test_value3"},"body":"dGVzdF9kYXRh"}`
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SRequestError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestInvalidRequest(t *testing.T) {
	t.Parallel()

	if _, err := LoadRequest([]byte{123}); err == nil {
		t.Error("success load invalid request bytes")
		return
	}

	bytesJoiner := joiner.NewBytesJoiner32([][]byte{
		{byte(123)},
		{byte(111)},
	})
	if _, err := LoadRequest(bytesJoiner); err == nil {
		t.Error("success load invalid request bytes joiner")
		return
	}

	if _, err := LoadRequest("123"); err == nil {
		t.Error("success load invalid request string")
		return
	}

	if _, err := LoadRequest(struct{}{}); err == nil {
		t.Error("success load invalid request type")
		return
	}
}

func TestRequest(t *testing.T) {
	t.Parallel()

	request := NewRequest(tcMethod, tcHost, tcPath).
		WithHead(tgHead).
		WithBody(tgBody)

	if request.GetHost() != tcHost {
		t.Error("host is not equals")
		return
	}

	if request.GetPath() != tcPath {
		t.Error("path is not equals")
		return
	}

	if request.GetMethod() != tcMethod {
		t.Error("method is not equals")
		return
	}

	for k, v := range request.GetHead() {
		v1, ok := tgHead[k]
		if !ok {
			t.Errorf("header undefined '%s'", k)
			return
		}
		if v != v1 {
			t.Errorf("header is invalid '%s'", v1)
			return
		}
	}

	if !bytes.Equal(request.GetBody(), tgBody) {
		t.Error("body is not equals")
		return
	}
}

func TestLoadRequest(t *testing.T) {
	t.Parallel()

	brequest := NewRequest(tcMethod, tcHost, tcPath).
		WithHead(tgHead).
		WithBody(tgBody).ToBytes()

	request1, err := LoadRequest(brequest)
	if err != nil {
		t.Error(err)
		return
	}

	request2, err := LoadRequest(tgBRequest)
	if err != nil {
		t.Error(err)
		return
	}

	reqStr := request2.ToString()
	if reqStr != tgBRequest {
		fmt.Println(reqStr)
		fmt.Println(tgBRequest)
		t.Error("string request is invalid")
		return
	}

	if request1.GetHost() != request2.GetHost() {
		t.Error("host is not equals")
		return
	}

	if request1.GetPath() != request2.GetPath() {
		t.Error("path is not equals")
		return
	}

	if request1.GetMethod() != request2.GetMethod() {
		t.Error("method is not equals")
		return
	}

	for k, v := range request1.GetHead() {
		v1, ok := request2.GetHead()[k]
		if !ok {
			t.Errorf("header undefined '%s'", k)
			return
		}
		if v != v1 {
			t.Errorf("header is invalid '%s'", v1)
			return
		}
	}

	if !bytes.Equal(request1.GetBody(), request2.GetBody()) {
		t.Error("body is not equals")
		return
	}
}
