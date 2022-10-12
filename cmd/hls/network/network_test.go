package network

import (
	"bytes"
	"testing"
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
	tgBRequest = []byte(`{
		"host": "test_host",
		"path": "test_path",
		"method": "test_method",
		"head": {
			"test_header1": "test_value1",
			"test_header2": "test_value2",
			"test_header3": "test_value3"
		},
		"body": "dGVzdF9kYXRh"
	}`)
)

func TestRequest(t *testing.T) {
	request := NewRequest(tcMethod, tcHost, tcPath).
		WithHead(tgHead).
		WithBody(tgBody)

	if request.Host() != tcHost {
		t.Error("host is not equals")
	}

	if request.Path() != tcPath {
		t.Error("path is not equals")
	}

	if request.Method() != tcMethod {
		t.Error("method is not equals")
	}

	for k, v := range request.Head() {
		v1, ok := tgHead[k]
		if !ok {
			t.Errorf("header undefined '%s'", k)
		}
		if v != v1 {
			t.Errorf("header is invalid '%s'", v1)
		}
	}

	if !bytes.Equal(request.Body(), tgBody) {
		t.Error("body is not equals")
	}
}

func TestLoadRequest(t *testing.T) {
	brequest := NewRequest(tcMethod, tcHost, tcPath).
		WithHead(tgHead).
		WithBody(tgBody).Bytes()

	request1 := LoadRequest(brequest)
	request2 := LoadRequest(tgBRequest)

	if request1.Host() != request2.Host() {
		t.Error("host is not equals")
	}

	if request1.Path() != request2.Path() {
		t.Error("path is not equals")
	}

	if request1.Method() != request2.Method() {
		t.Error("method is not equals")
	}

	for k, v := range request1.Head() {
		v1, ok := request2.Head()[k]
		if !ok {
			t.Errorf("header undefined '%s'", k)
		}
		if v != v1 {
			t.Errorf("header is invalid '%s'", v1)
		}
	}

	if !bytes.Equal(request1.Body(), request2.Body()) {
		t.Error("body is not equals")
	}
}
