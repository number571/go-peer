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

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SEncodingError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestHex(t *testing.T) {
	t.Parallel()

	data := HexEncode(tgNumInBytes)
	if !bytes.Equal(tgNumInBytes, HexDecode(data)) {
		t.Fatal("bytes not equals")
	}
	if dec := HexDecode("!@#"); dec != nil {
		t.Fatal("success decode invalid data")
	}
}

func TestBytes(t *testing.T) {
	t.Parallel()

	bnum := Uint64ToBytes(tgBytesInNum)
	if tgBytesInNum != BytesToUint64(bnum) {
		t.Fatal("numbers not equals")
	}
}

func TestSerializeJSON(t *testing.T) {
	t.Parallel()

	if string(SerializeJSON(tgMessage)) != tcJSON {
		t.Fatal("serialize string is invalid (non indent)")
	}

	res := new(tsMessage)

	if err := DeserializeJSON([]byte(tcJSON), res); err != nil {
		t.Fatal(err)
	}

	if res.FResult != "hello" || res.FReturn != 5 {
		t.Fatal("fields not equals")
	}

	if err := DeserializeJSON([]byte(`qwerty`), res); err == nil {
		t.Fatal("success deserialize invalid data")
	}
}

func TestSerializeYAML(t *testing.T) {
	t.Parallel()

	if string(SerializeYAML(tgMessage)) != tcYaml {
		t.Fatal("serialize string is invalid (non indent)")
	}

	res := new(tsMessage)

	if err := DeserializeYAML([]byte(tcYaml), res); err != nil {
		t.Fatal(err)
	}

	if res.FResult != "hello" || res.FReturn != 5 {
		t.Fatal("fields not equals")
	}

	if err := DeserializeYAML([]byte(`qwerty`), res); err == nil {
		t.Fatal("success deserialize invalid data")
	}
}
