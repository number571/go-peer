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
