package testutils

import (
	"time"

	"github.com/number571/go-peer/modules/client"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/friends"
	"github.com/number571/go-peer/modules/network"
	"github.com/number571/go-peer/modules/network/anonymity"
	"github.com/number571/go-peer/modules/network/conn"
	"github.com/number571/go-peer/modules/queue"
	"github.com/number571/go-peer/modules/storage/database"
	"github.com/number571/go-peer/settings/testutils"
)

func TestNewNode(pathDB string) anonymity.INode {
	msgSize := uint64(100 << 10)
	node := anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FTimeWait: 30 * time.Second,
		}),
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
