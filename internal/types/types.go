package types

import (
	"context"
	"sync"
)

type IServiceF func(context.Context, *sync.WaitGroup, chan<- error)
