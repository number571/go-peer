package filesystem

import (
	"os"
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
		return nil, err
	}
	return data, nil
}

func (p *sFile) Write(pData []byte) error {
	return os.WriteFile(p.fPath, pData, 0644)
}

func (p *sFile) IsExist() bool {
	_, err := os.Stat(p.fPath)
	return !os.IsNotExist(err)
}
