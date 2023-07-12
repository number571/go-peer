package conn

import (
	"net"
	"time"

	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
)

type IConn interface {
	types.ICloser

	GetSettings() ISettings
	GetSocket() net.Conn

	FetchPayload(pld payload.IPayload) (payload.IPayload, error)
	WritePayload(payload.IPayload) error
	ReadPayload() payload.IPayload
}

type ISettings interface {
	GetNetworkKey() string
	GetMessageSize() uint64
	GetLimitVoidSize() uint64
	GetFetchTimeWait() time.Duration
	GetWaitReadDeadline() time.Duration
	GetReadDeadline() time.Duration
	GetWriteDeadline() time.Duration
}
