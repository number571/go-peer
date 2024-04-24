package stream

import "io"

type IReadSeeker io.ReadSeeker

type IFileInfo interface {
	GetName() string
	GetHash() string
	GetSize() uint64
}
