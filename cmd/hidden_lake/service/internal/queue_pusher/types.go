package queue_pusher

import "github.com/number571/go-peer/pkg/queue_set"

type IQPWrapper interface {
	Get() queue_set.IQueuePusher
	Set(queue_set.IQueuePusher) IQPWrapper
}
