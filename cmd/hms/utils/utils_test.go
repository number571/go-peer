package utils

import (
	"os"
	"testing"

	"github.com/number571/go-peer/crypto"
)

type tsMessage struct {
	Result string `json:"result"`
	Return int    `json:"return"`
}

const (
	tcUtilsFile = "utils_test.txt"
	tcFileData  = `test text
for utils package
`
	tcJSON = `{
	"result": "hello",
	"return": 5
}`
)

var (
	tgMessage  = tsMessage{"hello", 5}
	tgRandFile = crypto.NewPRNG().String(20)
)

func TestFileIsExist(t *testing.T) {
	if FileIsExist(tgRandFile) {
		t.Errorf("file with name '%s' exists?", tgRandFile)
	}

	if !FileIsExist(tcUtilsFile) {
		t.Errorf("file with name '%s' does not exists?", tcUtilsFile)
	}
}

func TestReadFile(t *testing.T) {
	res := ReadFile(tgRandFile)
	if res != nil {
		t.Errorf("success read random file '%s'?", tgRandFile)
	}

	res = ReadFile(tcUtilsFile)
	if res == nil {
		t.Errorf("failed read file '%s'", tcUtilsFile)
	}

	if string(res) != tcFileData {
		t.Errorf("invalid read text from '%s'", tcUtilsFile)
	}
}

func TestWriteFile(t *testing.T) {
	defer os.Remove(tgRandFile)

	err := WriteFile(tgRandFile, []byte(tcFileData))
	if err != nil {
		t.Errorf("failed write to random file '%s'?", tgRandFile)
	}

	res := ReadFile(tgRandFile)
	if res == nil {
		t.Errorf("failed read random file '%s'", tgRandFile)
	}

	if string(res) != tcFileData {
		t.Errorf("invalid read text from '%s'", tgRandFile)
	}
}

func TestSerialize(t *testing.T) {
	if string(Serialize(tgMessage)) != tcJSON {
		t.Errorf("serialize string is invalid")
	}
}

func TestDeserialize(t *testing.T) {
	res := new(tsMessage)

	err := Deserialize([]byte(tcJSON), res)
	if err != nil {
		t.Errorf("deserialize failed")
	}

	if res.Result != "hello" || res.Return != 5 {
		t.Errorf("fields not equals")
	}
}
