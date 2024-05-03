package request

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
	tgBRequest = `{
		"host": "test_host",
		"path": "test_path",
		"method": "test_method",
		"head": {
			"test_header1": "test_value1",
			"test_header2": "test_value2",
			"test_header3": "test_value3"
		},
		"body": "dGVzdF9kYXRh"
	}`
)

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
