package database

import "io"

type IKVDatabase interface {
	io.Closer

	Set([]byte, []byte) error
	Get([]byte) ([]byte, error)
	Del([]byte) error
}
