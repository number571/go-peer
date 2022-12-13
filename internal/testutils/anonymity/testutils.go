package testutils

import (
	"time"

	"github.com/number571/go-peer/internal/testutils"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/friends"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/queue"
	"github.com/number571/go-peer/pkg/storage/database"
)

func TestNewNode(pathDB string) anonymity.INode {
	msgSize := uint64(100 << 10)
	node := anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FTimeWait: 30 * time.Second,
		}),
		logger.NewLogger(logger.NewSettings(&logger.SSettings{})),
		database.NewLevelDB(
			database.NewSettings(&database.SSettings{
				FHashing:   true,
				FCipherKey: []byte(testutils.TcKey1),
			}),
			pathDB,
		),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FCapacity:    (1 << 10),
				FMaxConnects: 10,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FMessageSize: msgSize,
					FTimeWait:    5 * time.Second,
				}),
			}),
		),
		queue.NewQueue(
			queue.NewSettings(&queue.SSettings{
				FCapacity:     10,
				FPullCapacity: 5,
				FDuration:     500 * time.Millisecond,
			}),
			client.NewClient(
				client.NewSettings(&client.SSettings{
					FWorkSize:    10,
					FMessageSize: msgSize,
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey),
			),
		),
		friends.NewF2F(),
	)
	return node
}
