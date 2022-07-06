package utils

type IFile interface {
	Read() ([]byte, error)
	Write([]byte) error
	IsExist() bool
}

type IInput interface {
	String() string
	Password() string
}

type ICloser interface {
	Close() error
}
