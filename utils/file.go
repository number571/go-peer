package utils

import (
	"io/ioutil"
	"os"
)

type IFile interface {
	Read() ([]byte, error)
	Write([]byte) error
	IsExist() bool
}

var (
	_ IFile = &sFile{}
)

type sFile struct {
	fPath string
}

func OpenFile(path string) IFile {
	return &sFile{
		fPath: path,
	}
}

func (file *sFile) Read() ([]byte, error) {
	data, err := ioutil.ReadFile(file.fPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (file *sFile) Write(data []byte) error {
	return ioutil.WriteFile(file.fPath, data, 0644)
}

func (file *sFile) IsExist() bool {
	_, err := os.Stat(file.fPath)
	return !os.IsNotExist(err)
}
