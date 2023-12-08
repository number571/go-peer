package encoding

import (
	"bytes"
	"testing"
)

type tsMessage struct {
	FResult string `yaml:"result" json:"result"`
	FReturn int    `yaml:"return" json:"return"`
}

const (
	tgBytesInNum = uint64(0xABCDEF0123456789)
	tcJSON       = `{"result":"hello","return":5}`
	tcYaml       = `result: hello
return: 5
`
)

var (
	tgNumInBytes = []byte{0, 1, 2, 3, 4, 5, 6, 7, 255, 254, 253, 252, 251, 250, 128, 127, 126, 125}
	tgMessage    = tsMessage{"hello", 5}
)

func TestHex(t *testing.T) {
	t.Parallel()

	data := HexEncode(tgNumInBytes)
	if !bytes.Equal(tgNumInBytes, HexDecode(data)) {
		t.Error("bytes not equals")
		return
	}
	if dec := HexDecode("!@#"); dec != nil {
		t.Error("success decode invalid data")
		return
	}
}

func TestBytes(t *testing.T) {
	t.Parallel()

	bnum := Uint64ToBytes(tgBytesInNum)
	if tgBytesInNum != BytesToUint64(bnum) {
		t.Error("numbers not equals")
		return
	}
}

func TestSerializeJSON(t *testing.T) {
	t.Parallel()

	if string(SerializeJSON(tgMessage)) != tcJSON {
		t.Error("serialize string is invalid (non indent)")
		return
	}

	res := new(tsMessage)

	if err := DeserializeJSON([]byte(tcJSON), res); err != nil {
		t.Error(err)
		return
	}

	if res.FResult != "hello" || res.FReturn != 5 {
		t.Error("fields not equals")
		return
	}

	if err := DeserializeJSON([]byte(`qwerty`), res); err == nil {
		t.Error("success deserialize invalid data")
		return
	}
}

func TestSerializeYAML(t *testing.T) {
	t.Parallel()

	if string(SerializeYAML(tgMessage)) != tcYaml {
		t.Error("serialize string is invalid (non indent)")
		return
	}

	res := new(tsMessage)

	if err := DeserializeYAML([]byte(tcYaml), res); err != nil {
		t.Error(err)
		return
	}

	if res.FResult != "hello" || res.FReturn != 5 {
		t.Error("fields not equals")
		return
	}

	if err := DeserializeYAML([]byte(`qwerty`), res); err == nil {
		t.Error("success deserialize invalid data")
		return
	}
}
