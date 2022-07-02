package utils

type ICloser interface {
	Close() error
}

type IFile interface {
	Read() ([]byte, error)
	Write([]byte) error
	IsExist() bool
}
