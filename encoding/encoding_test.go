package encoding

import (
	"bytes"
	"testing"
)

type tsMessage struct {
	Result string `json:"result"`
	Return int    `json:"return"`
}

const (
	tgBytesInNum = uint64(0xABCDEF0123456789)
	tcJSON       = `{
	"result": "hello",
	"return": 5
}`
)

var (
	tgNumInBytes = []byte{0, 1, 2, 3, 4, 5, 6, 7, 255, 254, 253, 252, 251, 250, 128, 127, 126, 125}
	tgMessage    = tsMessage{"hello", 5}
)

func TestBase64(t *testing.T) {
	data := Base64Encode(tgNumInBytes)
	if !bytes.Equal(tgNumInBytes, Base64Decode(data)) {
		t.Error("bytes not equals")
	}
}

func TestBytes(t *testing.T) {
	bnum := Uint64ToBytes(tgBytesInNum)
	if tgBytesInNum != BytesToUint64(bnum) {
		t.Error("numbers not equals")
	}
}

func TestSerialize(t *testing.T) {
	if string(Serialize(tgMessage)) != tcJSON {
		t.Error("serialize string is invalid")
	}

	res := new(tsMessage)

	err := Deserialize([]byte(tcJSON), res)
	if err != nil {
		t.Error("deserialize failed")
	}

	if res.Result != "hello" || res.Return != 5 {
		t.Error("fields not equals")
	}
}
