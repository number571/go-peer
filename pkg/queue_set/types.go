package queue_set

type IQueueSet interface {
	GetSettings() ISettings

	GetIndex() uint64
	GetKey(i uint64) ([]byte, bool)

	IQueuePusher
	IQueueLoader
}

type IQueuePusher interface {
	Push([]byte, []byte) bool
}

type IQueueLoader interface {
	Load([]byte) ([]byte, bool)
}

type ISettings interface {
	GetCapacity() uint64
}
