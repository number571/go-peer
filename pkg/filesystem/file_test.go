package filesystem

import (
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
)

const (
	tcRandFile1 = "somerandfile1"
	tcRandFile2 = "somerandfile2"
	tcRandFile3 = "somerandfile3"
)

const (
	tcUtilsFile = "file_test.txt"
	tcFileData  = `test text
for filesystem package
`
)

func TestFileIsExist(t *testing.T) {
	if OpenFile(tcRandFile1).IsExist() {
		t.Errorf("file with name '%s' exists?", tcRandFile1)
	}

	if !OpenFile(tcUtilsFile).IsExist() {
		t.Errorf("file with name '%s' does not exists?", tcUtilsFile)
	}
}

func TestReadFile(t *testing.T) {
	_, err := OpenFile(tcRandFile2).Read()
	if err == nil {
		t.Errorf("success read random file '%s'?", tcRandFile2)
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
	defer os.Remove(tcRandFile3)

	err := OpenFile(tcRandFile3).Write([]byte(tcFileData))
	if err != nil {
		t.Errorf("failed write to random file '%s'?", tcRandFile3)
	}

	res, err := OpenFile(tcRandFile3).Read()
	if err != nil {
		t.Errorf("failed read random file '%s'", tcRandFile3)
	}

	if string(res) != tcFileData {
		t.Errorf("invalid read text from '%s'", tcRandFile3)
	}

	prng := random.NewStdPRNG()
	randInvalidPath := prng.GetString(32) + "/" + prng.GetString(32) + "/" + prng.GetString(32)
	if err := OpenFile(randInvalidPath).Write([]byte("hello, world!")); err == nil {
		t.Error("success write bytes to invalid path")
		return
	}
}
