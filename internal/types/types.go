package types

import (
	"context"
	"sync"
)

type IService func(context.Context, *sync.WaitGroup, chan<- error)
