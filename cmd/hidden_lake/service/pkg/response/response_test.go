package response

import (
	"bytes"
	"testing"
)

const (
	tcResponse = `{"code":200,"head":{"key1":"value1","key2":"value2","key3":"value3"},"body":"aGVsbG8sIHdvcmxkIQ=="}`
	tcBody     = "hello, world!"
)

var (
	tgHead = map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
)

func TestResponse(t *testing.T) {
	t.Parallel()

	resp := NewResponse(200).
		WithHead(tgHead).
		WithBody([]byte(tcBody))

	resp1, err := LoadResponse(resp.ToBytes())
	if err != nil {
		t.Error(err)
		return
	}

	respStr := resp.ToString()
	if respStr != tcResponse {
		t.Error("string response is invalid")
		return
	}

	resp2, err := LoadResponse(respStr)
	if err != nil {
		t.Error(err)
		return
	}

	testResponse(t, resp)
	testResponse(t, resp1)
	testResponse(t, resp2)
}

func testResponse(t *testing.T, resp IResponse) {
	if resp.GetCode() != 200 {
		t.Error("resp code is invalid")
		return
	}
	if !bytes.Equal(resp.GetBody(), []byte(tcBody)) {
		t.Error("resp body is invalid")
		return
	}
	if len(resp.GetHead()) != 3 {
		t.Error("resp head size is invalid")
		return
	}

	for k, v := range resp.GetHead() {
		v1, ok := tgHead[k]
		if !ok {
			t.Error("undefined value in orig head")
			return
		}
		if v1 != v {
			t.Error("resp head value is invalid")
			return
		}
	}
}
