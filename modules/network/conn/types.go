package conn

import (
	"net"
	"time"

	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/payload"
)

type ISettings interface {
	GetNetworkKey() string
	GetMessageSize() uint64
	GetTimeWait() time.Duration
}

type IConn interface {
	Settings() ISettings
	modules.ICloser

	Socket() net.Conn
	Request(pld payload.IPayload) (payload.IPayload, error)

	Write(payload.IPayload) error
	Read() payload.IPayload
}
