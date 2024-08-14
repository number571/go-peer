package friends

type IListKeys interface {
	GetKeys() [][]byte
	AddKey([]byte)
	DelKey([]byte)
}
