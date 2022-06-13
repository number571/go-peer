package utils

import (
	"io/ioutil"
	"os"
)

var (
	_ IFile = &sFile{}
)

type sFile struct {
	path string
}

func NewFile(path string) IFile {
	return &sFile{
		path: path,
	}
}

func (file *sFile) Read() ([]byte, error) {
	data, err := ioutil.ReadFile(file.path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (file *sFile) Write(data []byte) error {
	return ioutil.WriteFile(file.path, data, 0644)
}

func (file *sFile) IsExist() bool {
	_, err := os.Stat(file.path)
	return !os.IsNotExist(err)
}
