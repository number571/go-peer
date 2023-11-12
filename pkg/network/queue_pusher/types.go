package queue_pusher

type IQueuePusher interface {
	Push([]byte) bool
}

type ISettings interface {
	GetCapacity() uint64
}
