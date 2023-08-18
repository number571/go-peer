package keybuilder

type IKeyBuilder interface {
	Build([]byte) []byte
}
