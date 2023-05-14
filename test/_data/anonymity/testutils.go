package testutils

import (
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/database"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	TCMessageSize = uint64(100 << 10)
	TCWorkSize    = 10
)

func TestNewNode(dbPath, addr string) anonymity.INode {
	db, err := database.NewSQLiteDB(
		database.NewSettings(&database.SSettings{
			FPath:      dbPath,
			FHashing:   true,
			FCipherKey: []byte("CIPHER"),
		}),
	)
	if err != nil {
		return nil
	}
	node := anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FTimeWait: 30 * time.Second,
		}),
		logger.NewLogger(logger.NewSettings(&logger.SSettings{})),
		anonymity.NewWrapperDB().Set(db),
		TestNewNetworkNode(addr),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{
				FCapacity:     10,
				FPullCapacity: 5,
				FDuration:     500 * time.Millisecond,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FWorkSize:    TCWorkSize,
					FMessageSize: TCMessageSize,
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey),
			),
		),
		asymmetric.NewListPubKeys(),
	)
	return node
}

func TestNewNetworkNode(addr string) network.INode {
	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:     addr,
			FCapacity:    (1 << 10),
			FMaxConnects: 10,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSize: TCMessageSize,
				FTimeWait:    5 * time.Second,
			}),
		}),
	)
}
