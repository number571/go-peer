package keybuilder

type IKeyBuilder interface {
	Build(string, uint64) []byte
}
