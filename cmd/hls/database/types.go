package database

type IKeyValueDB interface {
	Push([]byte) error
	Exist([]byte) bool
	Close() error
}
