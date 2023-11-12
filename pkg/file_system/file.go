package file_system

import (
	"os"

	"github.com/number571/go-peer/pkg/errors"
)

var (
	_ IFile = &sFile{}
)

type sFile struct {
	fPath string
}

func OpenFile(pPath string) IFile {
	return &sFile{
		fPath: pPath,
	}
}

func (p *sFile) Read() ([]byte, error) {
	data, err := os.ReadFile(p.fPath)
	if err != nil {
		return nil, errors.WrapError(err, "read file")
	}
	return data, nil
}

func (p *sFile) Write(pData []byte) error {
	if err := os.WriteFile(p.fPath, pData, 0644); err != nil {
		return errors.WrapError(err, "write file")
	}
	return nil
}

func (p *sFile) IsExist() bool {
	_, err := os.Stat(p.fPath)
	return !os.IsNotExist(err)
}
