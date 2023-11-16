package queue_set

type IQueueSet interface {
	GetSettings() ISettings

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