package keybuilder

type IKeyBuilder interface {
	Build(string) []byte
}
