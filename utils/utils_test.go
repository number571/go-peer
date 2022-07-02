package utils

import (
	"os"
	"testing"

	"github.com/number571/go-peer/crypto/random"
)

const (
	tcUtilsFile = "utils_test.txt"
	tcFileData  = `test text
for utils package
`
)

var (
	tgRandFile = random.NewStdPRNG().String(20)
)

func TestFileIsExist(t *testing.T) {
	if OpenFile(tgRandFile).IsExist() {
		t.Errorf("file with name '%s' exists?", tgRandFile)
	}

	if !OpenFile(tcUtilsFile).IsExist() {
		t.Errorf("file with name '%s' does not exists?", tcUtilsFile)
	}
}

func TestReadFile(t *testing.T) {
	_, err := OpenFile(tgRandFile).Read()
	if err == nil {
		t.Errorf("success read random file '%s'?", tgRandFile)
	}

	res, err := OpenFile(tcUtilsFile).Read()
	if err != nil {
		t.Errorf("failed read file '%s'", tcUtilsFile)
	}

	if string(res) != tcFileData {
		t.Errorf("invalid read text from '%s'", tcUtilsFile)
	}
}

func TestWriteFile(t *testing.T) {
	defer os.Remove(tgRandFile)

	err := OpenFile(tgRandFile).Write([]byte(tcFileData))
	if err != nil {
		t.Errorf("failed write to random file '%s'?", tgRandFile)
	}

	res, err := OpenFile(tgRandFile).Read()
	if err != nil {
		t.Errorf("failed read random file '%s'", tgRandFile)
	}

	if string(res) != tcFileData {
		t.Errorf("invalid read text from '%s'", tgRandFile)
	}
}
